package middleware

import (
	"Test/config"
	"Test/pkg/jwt"
	"github.com/DmytroBuzhylov/ratelimiter"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"
)

type AuthMiddleware struct {
	config      config.JWTConfig
	rateLimiter *ratelimiter.RedisTokenBucket
}

func NewAuthMiddleware(config config.JWTConfig, rateLimiter *ratelimiter.RedisTokenBucket) *AuthMiddleware {
	return &AuthMiddleware{config: config, rateLimiter: rateLimiter}
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

var requestBucket = make(map[string]int)

func (m *AuthMiddleware) RateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {

		ok, err := m.rateLimiter.AllowIP(c.Request.Context(), c.ClientIP())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "redis error",
			})
			log.Warn().Err(err).
				Str("middleware", "AuthMiddleware").
				Str("method", "RateLimiter").
				Msg("redis error")
			c.Abort()
			return
		}
		if !ok {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "to many req",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
