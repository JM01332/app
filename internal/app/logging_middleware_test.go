package app

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestRequestLoggingMiddlewareLogsCompletedRequest(t *testing.T) {
	logger, logs := observedLogger()
	router := gin.New()
	router.Use(requestLoggingMiddleware(logger), recoveryMiddleware(logger))
	router.GET("/test", func(context *gin.Context) {
		context.Status(http.StatusNoContent)
	})

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/test", nil)
	router.ServeHTTP(response, request)

	entries := logs.FilterMessage("HTTP request").All()
	if len(entries) != 1 {
		t.Fatalf("HTTP request log count = %d, want 1", len(entries))
	}
	if entries[0].Level != zapcore.InfoLevel {
		t.Errorf("log level = %v, want info", entries[0].Level)
	}

	fields := entries[0].ContextMap()
	if fields["method"] != http.MethodGet {
		t.Errorf("logged method = %v, want GET", fields["method"])
	}
	if fields["path"] != "/test" {
		t.Errorf("logged path = %v, want /test", fields["path"])
	}
	if fields["status"] != int64(http.StatusNoContent) {
		t.Errorf("logged status = %v, want %d", fields["status"], http.StatusNoContent)
	}
}

func TestRequestLoggingMiddlewareLogsServerError(t *testing.T) {
	logger, logs := observedLogger()
	router := gin.New()
	router.Use(requestLoggingMiddleware(logger), recoveryMiddleware(logger))
	router.GET("/failure", func(context *gin.Context) {
		_ = context.Error(errors.New("database unavailable"))
		context.Status(http.StatusInternalServerError)
	})

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/failure", nil)
	router.ServeHTTP(response, request)

	entries := logs.FilterMessage("HTTP request").All()
	if len(entries) != 1 {
		t.Fatalf("HTTP request log count = %d, want 1", len(entries))
	}
	if entries[0].Level != zapcore.ErrorLevel {
		t.Errorf("log level = %v, want error", entries[0].Level)
	}
	if _, ok := entries[0].ContextMap()["errors"]; !ok {
		t.Error("HTTP error log does not contain errors field")
	}
}

func TestRecoveryMiddlewareLogsPanicAndReturnsControlledError(t *testing.T) {
	logger, logs := observedLogger()
	router := gin.New()
	router.Use(requestLoggingMiddleware(logger), recoveryMiddleware(logger))
	router.GET("/panic", func(context *gin.Context) {
		panic("unexpected failure")
	})

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/panic", nil)
	router.ServeHTTP(response, request)

	if response.Code != http.StatusInternalServerError {
		t.Fatalf("status code = %d, want %d", response.Code, http.StatusInternalServerError)
	}
	if logs.FilterMessage("panic recovered").Len() != 1 {
		t.Fatal("panic recovery log count != 1")
	}
	if logs.FilterMessage("HTTP request").Len() != 1 {
		t.Fatal("HTTP request log count != 1")
	}
}

func observedLogger() (*zap.Logger, *observer.ObservedLogs) {
	core, logs := observer.New(zapcore.DebugLevel)
	return zap.New(core), logs
}
