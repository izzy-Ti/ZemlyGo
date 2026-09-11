package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/izzy-Ti/ZemlyGo/internals/constants"
	"github.com/izzy-Ti/ZemlyGo/internals/dto"
)

func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get("user_role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, dto.ErrorAPIResponse("Role not found in token context"))
			return
		}

		userRole, ok := roleVal.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, dto.ErrorAPIResponse("Invalid role in context"))
			return
		}

		// Admin role always has access
		if userRole == constants.RoleAdmin {
			c.Next()
			return
		}

		for _, allowed := range allowedRoles {
			if userRole == allowed {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, dto.ErrorAPIResponse("Forbidden: requires one of roles "+formatRoles(allowedRoles)))
	}
}

func formatRoles(roles []string) string {
	var s string
	for i, r := range roles {
		if i > 0 {
			s += ", "
		}
		s += r
	}
	return "[" + s + "]"
}
