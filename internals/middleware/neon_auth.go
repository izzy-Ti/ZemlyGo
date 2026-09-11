package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/izzy-Ti/ZemlyGo/internals/domain"
	"github.com/izzy-Ti/ZemlyGo/internals/dto"
	"github.com/izzy-Ti/ZemlyGo/internals/infrastructure/neon"
	"github.com/izzy-Ti/ZemlyGo/internals/repository/interfaces"
)

func NeonAuthMiddleware(neonClient *neon.Client, userRepo interfaces.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorAPIResponse("Authorization header is required"))
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader || token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorAPIResponse("Invalid authorization token format"))
			return
		}

		claims, err := neonClient.ValidateToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorAPIResponse("Invalid or expired authentication token: "+err.Error()))
			return
		}

		// Look up user by Neon Auth ID
		user, err := userRepo.GetByAuthID(claims.AuthID)
		if err != nil || user == nil {
			// Check by email
			if claims.Email != "" {
				user, _ = userRepo.GetByEmail(claims.Email)
			}
		}

		if user == nil {
			// Auto-register user authenticated via Neon Auth on first API call
			role := claims.Role
			if role == "" {
				role = "rider"
			}
			name := claims.Name
			if name == "" {
				name = claims.Email
			}

			user = &domain.Users{
				AuthID:    claims.AuthID,
				Email:     claims.Email,
				Name:      name,
				Phone:     "+0000000000",
				Role:      role,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			_ = userRepo.Create(user)
		}

		c.Set("user_id", user.ID)
		c.Set("auth_id", user.AuthID)
		c.Set("supabase_user_id", user.AuthID) // backward compatibility
		c.Set("user_email", user.Email)
		c.Set("user_role", user.Role)
		c.Set("current_user", user)

		c.Next()
	}
}
