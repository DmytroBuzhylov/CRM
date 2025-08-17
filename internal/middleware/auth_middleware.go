package middleware

import (
	"Test/config"
	"Test/pkg/jwt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"strings"
)

type AuthMiddleware struct {
	config config.JWTConfig
}

func NewAuthMiddleware(config config.JWTConfig) *AuthMiddleware {
	return &AuthMiddleware{config: config}
}

func (m *AuthMiddleware) JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "Токен отсутствует"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(401, gin.H{"error": "Неверный формат токена"})
			c.Abort()
			return
		}
		claims, err := jwt.VerifyToken(parts[1], m.config.JWTAccessSecret)
		log.Info().Err(err).Send()
		if err != nil {
			c.JSON(401, gin.H{"error": "Неверный или истекший токен"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		c.Set("organization_id", claims.OrganizationID)

		c.Next()
	}
}
