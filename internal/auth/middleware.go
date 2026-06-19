package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Middleware protects routes with a bearer token check.
func Middleware(verifier TokenVerifier) gin.HandlerFunc {
	return func(context *gin.Context) {
		token, ok := bearerToken(context.GetHeader("Authorization"))
		if !ok {
			abortUnauthorized(context)
			return
		}

		if err := verifier.Verify(context.Request.Context(), token); err != nil {
			_ = context.Error(err)
			abortUnauthorized(context)
			return
		}

		context.Next()
	}
}

func bearerToken(header string) (string, bool) {
	scheme, token, ok := strings.Cut(strings.TrimSpace(header), " ")
	if !ok || !strings.EqualFold(scheme, "Bearer") {
		return "", false
	}

	token = strings.TrimSpace(token)
	if token == "" {
		return "", false
	}

	return token, true
}

func abortUnauthorized(context *gin.Context) {
	context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"error": gin.H{
			"code":    "unauthorized",
			"message": "Bearer token is missing or invalid",
			"fields":  []any{},
		},
	})
}
