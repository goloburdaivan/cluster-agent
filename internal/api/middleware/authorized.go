package middleware

import (
	"cluster-agent/internal/auth"
	"cluster-agent/internal/auth/permissions"
	"cluster-agent/internal/config"
	"github.com/gin-gonic/gin"
	"net/http"
	"slices"
	"strings"
)

type AuthorizedMiddleware struct {
	cfg *config.Config
}

func NewAuthorizedMiddleware(cfg *config.Config) *AuthorizedMiddleware {
	return &AuthorizedMiddleware{
		cfg: cfg,
	}
}

func (m *AuthorizedMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenStr string

		authHeader := c.GetHeader("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenStr = authHeader[7:]
		}

		if tokenStr == "" {
			tokenStr = c.Query("token")
		}

		if tokenStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing auth token"})
			return
		}

		claims := &auth.UserClaims{}
		parsedToken, err := auth.ParseToken(tokenStr, claims, m.cfg.JWTPublicKey)

		if err != nil || !parsedToken.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		c.Set("claims", claims)
		c.Next()
	}
}

func (m *AuthorizedMiddleware) HasPermission(permission permissions.Permission) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := GetUserClaims(c)

		if claims == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing jwt token claims"})
			return
		}

		if !slices.Contains(claims.Permissions, permission) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "invalid permission"})
			return
		}

		c.Next()
	}
}

func GetUserClaims(c *gin.Context) *auth.UserClaims {
	if val, exists := c.Get("claims"); exists {
		if claims, ok := val.(*auth.UserClaims); ok {
			return claims
		}
	}

	return nil
}
