package middleware

import (
	"Test/internal/feature/permission/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type PermissionMiddleware struct {
	uc usecase.PermissionUseCase
}

func NewPermissionMiddleware(uc usecase.PermissionUseCase) *PermissionMiddleware {
	return &PermissionMiddleware{
		uc: uc,
	}
}

func (m *PermissionMiddleware) CheckPermission(permissionName string) gin.HandlerFunc {
	return func(c *gin.Context) {

		userIDStr, ok := c.Get("user_id")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "", //
			})
			c.Abort()
			return
		}
		userID, ok := userIDStr.(uuid.UUID)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "", //
			})
			c.Abort()
			return
		}

		organizationIDStr, ok := c.Get("organization_id")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "", //
			})
			c.Abort()
			return
		}
		organizationID, ok := organizationIDStr.(uuid.UUID)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "", //
			})
			c.Abort()
			return
		}

		ok = m.uc.HasPermission(c.Request.Context(), userID, organizationID, permissionName)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "", //
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
