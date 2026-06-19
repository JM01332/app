package security

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const claimsContextKey = "security.claims"

// Authenticate verifies the Bearer token and stores its claims in the request context.
func Authenticate(verifier Verifier) gin.HandlerFunc {
	return func(context *gin.Context) {
		rawToken, ok := bearerToken(context.GetHeader("Authorization"))
		if !ok {
			abortUnauthorized(context)
			return
		}

		claims, err := verifier.Verify(context.Request.Context(), rawToken)
		if err != nil {
			abortUnauthorized(context)
			return
		}

		context.Set(claimsContextKey, claims)
		context.Next()
	}
}

// RequireAnyRole allows requests with at least one required Keycloak client role.
func RequireAnyRole(clientID string, requiredRoles ...Role) gin.HandlerFunc {
	return func(context *gin.Context) {
		value, exists := context.Get(claimsContextKey)
		claims, ok := value.(Claims)
		if !exists || !ok {
			abortUnauthorized(context)
			return
		}

		if !claims.HasAnyRole(clientID, requiredRoles...) {
			context.AbortWithStatusJSON(http.StatusForbidden, securityErrorResponse(
				"forbidden",
				"Insufficient permissions",
			))
			return
		}

		context.Next()
	}
}

func bearerToken(authorizationHeader string) (string, bool) {
	parts := strings.Fields(authorizationHeader)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || strings.TrimSpace(parts[1]) == "" {
		return "", false
	}

	return parts[1], true
}

func abortUnauthorized(context *gin.Context) {
	context.Header("WWW-Authenticate", "Bearer")
	context.AbortWithStatusJSON(http.StatusUnauthorized, securityErrorResponse(
		"unauthorized",
		"Missing or invalid access token",
	))
}

func securityErrorResponse(code string, message string) gin.H {
	return gin.H{
		"error": gin.H{
			"code":    code,
			"message": message,
			"fields":  []any{},
		},
	}
}
