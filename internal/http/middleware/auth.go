package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yourorg/shadowchat/backend/internal/service/auth"
)

func Auth(jwtSecret string) gin.HandlerFunc {
	authSvc := &auth.AuthService{} // In production, inject properly
	
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization format"})
			c.Abort()
			return
		}

		// For now, decode JWT manually - in production use proper middleware
		claims, err := authSvc.ValidateJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		// Check if token is blacklisted
		blacklisted, _ := authSvc.IsTokenBlacklisted(c.Request.Context(), claims.JTI)
		if blacklisted {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token revoked"})
			c.Abort()
			return
		}

		c.Set("userId", claims.UserID)
		c.Set("jti", claims.JTI)
		c.Next()
	}
}
