package middlewares

import (
	"Savannah_Screening_Test/clients"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
)

// RequireAuth is a middleware that verifies the JWT and checks for optional roles
func RequireAuth(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid Authorization header"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Make sure the signing method is RSA (used by Keycloak by default)
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			// Get public key from JWKS endpoint
			// You should cache this or use a lib like `lestrrat-go/jwx` for JWKS management
			key, err := clients.GetKeycloakPublicKey()
			if err != nil {
				return nil, err
			}
			return key, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		// Extract roles from claims
		roles := extractRolesFromClaims(claims)
		if !hasRequiredRoles(roles, requiredRoles) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: insufficient role permissions"})
			c.Abort()
			return
		}

		// Optionally pass user info to context
		c.Set("user", claims)
		c.Next()
	}
}

func hasRequiredRoles(userRoles, requiredRoles []string) bool {
	if len(requiredRoles) == 0 {
		return true // no role restriction
	}
	for _, required := range requiredRoles {
		for _, role := range userRoles {
			if role == required {
				return true
			}
		}
	}
	return false
}

func extractRolesFromClaims(claims jwt.MapClaims) []string {
	roles := []string{}

	if realmAccess, ok := claims["realm_access"].(map[string]interface{}); ok {
		if rolesClaim, ok := realmAccess["roles"].([]interface{}); ok {
			for _, r := range rolesClaim {
				if roleStr, ok := r.(string); ok {
					roles = append(roles, roleStr)
				}
			}
		}
	}

	return roles
}
