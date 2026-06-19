package security

import "strings"

// Role is an application role derived from Keycloak client roles.
type Role string

const (
	RoleUser  Role = "USER"
	RoleAdmin Role = "ADMIN"
)

// Claims contains the access-token claims used by this application.
type Claims struct {
	Subject           string                  `json:"sub"`
	PreferredUsername string                  `json:"preferred_username"`
	AuthorizedParty   string                  `json:"azp"`
	ResourceAccess    map[string]ClientAccess `json:"resource_access"`
}

// ClientAccess contains roles assigned for one Keycloak client.
type ClientAccess struct {
	Roles []string `json:"roles"`
}

// HasAnyRole reports whether the configured client has one required role.
func (claims Claims) HasAnyRole(clientID string, requiredRoles ...Role) bool {
	clientAccess, ok := claims.ResourceAccess[clientID]
	if !ok {
		return false
	}

	for _, tokenRole := range clientAccess.Roles {
		normalizedRole := Role(strings.ToUpper(tokenRole))
		for _, requiredRole := range requiredRoles {
			if normalizedRole == requiredRole {
				return true
			}
		}
	}

	return false
}
