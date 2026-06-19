package router

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/JM01332/app/internal/carrier/model"
	carrierservice "github.com/JM01332/app/internal/carrier/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

const (
	errorCodeInvalidJSON       = "invalid_json"
	errorCodeValidationFailed  = "validation_failed"
	errorCodeCarrierNameExists = "carrier_name_exists"
	errorCodeCarrierNotFound   = "carrier_not_found"
	errorCodeInternal          = "internal_error"
)

type CarrierService interface {
	List(ctx context.Context) ([]model.Carrier, error)
	GetByID(ctx context.Context, id int64) (*model.Carrier, error)
	Create(ctx context.Context, input carrierservice.CreateCarrierInput) (*model.Carrier, error)
}

type Handler struct {
	service  CarrierService
	validate *validator.Validate
}

func NewHandler(service CarrierService) *Handler {
	return &Handler{
		service:  service,
		validate: validator.New(),
	}
}

func RegisterRoutes(router gin.IRouter, service CarrierService) {
	handler := NewHandler(service)

	router.GET("/carriers", handler.List)
	router.GET("/carriers/:id", handler.GetByID)
	router.POST("/carriers", handler.Create)
}

func (handler *Handler) List(context *gin.Context) {
	carriers, err := handler.service.List(context.Request.Context())
	if err != nil {
		context.JSON(http.StatusInternalServerError, newErrorResponse(errorCodeInternal, "Internal server error", nil))
		return
	}

	context.JSON(http.StatusOK, mapCarrierResponses(carriers))
}

func (handler *Handler) GetByID(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		context.JSON(http.StatusBadRequest, newErrorResponse(errorCodeValidationFailed, "Carrier ID must be a positive integer", []FieldError{
			{Field: "id", Message: "id must be a positive integer"},
		}))
		return
	}

	carrier, err := handler.service.GetByID(context.Request.Context(), id)
	if errors.Is(err, carrierservice.ErrCarrierNotFound) {
		context.JSON(http.StatusNotFound, newErrorResponse(errorCodeCarrierNotFound, "Carrier not found", nil))
		return
	}
	if err != nil {
		context.JSON(http.StatusInternalServerError, newErrorResponse(errorCodeInternal, "Internal server error", nil))
		return
	}

	context.JSON(http.StatusOK, mapCarrierResponse(*carrier))
}

func (handler *Handler) Create(context *gin.Context) {
	var request CreateCarrierRequest
	if err := decodeStrictJSON(context.Request.Body, &request); err != nil {
		context.JSON(http.StatusBadRequest, newErrorResponse(errorCodeInvalidJSON, "Request body must be valid JSON", nil))
		return
	}

	if err := handler.validate.Struct(request); err != nil {
		context.JSON(http.StatusUnprocessableEntity, validationErrorResponse(err))
		return
	}

	input := mapCreateCarrierInput(request)
	carrier, err := handler.service.Create(context.Request.Context(), input)
	if errors.Is(err, carrierservice.ErrCarrierNameExists) {
		context.JSON(http.StatusConflict, newErrorResponse(errorCodeCarrierNameExists, "Carrier name already exists", nil))
		return
	}
	if err != nil {
		context.JSON(http.StatusInternalServerError, newErrorResponse(errorCodeInternal, "Internal server error", nil))
		return
	}

	context.Header("Location", "/api/carriers/"+strconvFormatInt64(carrier.ID))
	context.Status(http.StatusCreated)
}

func decodeStrictJSON(body io.Reader, target any) error {
	decoder := json.NewDecoder(body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(target); err != nil {
		return err
	}

	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return errors.New("request body must contain one JSON object")
	}

	return nil
}

func validationErrorResponse(err error) ErrorResponse {
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return newErrorResponse(errorCodeValidationFailed, "Request validation failed", nil)
	}

	fields := make([]FieldError, 0, len(validationErrors))
	for _, fieldError := range validationErrors {
		fields = append(fields, FieldError{
			Field:   jsonFieldName(fieldError),
			Message: validationMessage(fieldError),
		})
	}

	return newErrorResponse(errorCodeValidationFailed, "Request validation failed", fields)
}

func newErrorResponse(code string, message string, fields []FieldError) ErrorResponse {
	if fields == nil {
		fields = []FieldError{}
	}

	return ErrorResponse{
		Error: APIError{
			Code:    code,
			Message: message,
			Fields:  fields,
		},
	}
}

func jsonFieldName(fieldError validator.FieldError) string {
	fieldNames := map[string]string{
		"Name":          "name",
		"Nation":        "nation",
		"CarrierType":   "carrierType",
		"CommandCenter": "commandCenter",
		"CodeName":      "codeName",
		"SecurityLevel": "securityLevel",
		"Aircrafts":     "aircrafts",
		"Model":         "model",
		"Manufacturer":  "manufacturer",
	}

	fieldName, ok := fieldNames[fieldError.Field()]
	if !ok {
		return fieldError.Field()
	}

	namespace := fieldError.Namespace()
	if strings.Contains(namespace, "CommandCenter.") {
		return "commandCenter." + fieldName
	}
	if strings.Contains(namespace, "Aircrafts[") {
		indexStart := strings.Index(namespace, "Aircrafts[")
		indexEnd := strings.Index(namespace[indexStart:], "]")
		if indexStart >= 0 && indexEnd >= 0 {
			index := namespace[indexStart+len("Aircrafts[") : indexStart+indexEnd]
			return "aircrafts[" + index + "]." + fieldName
		}
	}

	return fieldName
}

func validationMessage(fieldError validator.FieldError) string {
	field := jsonFieldName(fieldError)

	switch fieldError.Tag() {
	case "required":
		return field + " is required"
	case "min":
		return field + " is too short"
	case "max":
		return field + " is too long"
	case "oneof":
		return field + " has an unsupported value"
	default:
		return field + " is invalid"
	}
}

func strconvFormatInt64(value int64) string {
	return strconv.FormatInt(value, 10)
}
