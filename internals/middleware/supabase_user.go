package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/izzy-Ti/ZemlyGo/configs"
	"github.com/izzy-Ti/ZemlyGo/internals/infrastructure/neon"
	"github.com/izzy-Ti/ZemlyGo/internals/repository/interfaces"
	"github.com/nedpals/supabase-go"
)

// AuthMiddleware provides backwards compatibility and delegates to NeonAuthMiddleware
func AuthMiddleware(supabaseClient *supabase.Client, userRepo interfaces.UserRepository, jwtSecret string) gin.HandlerFunc {
	neonClient := neon.NewClient(&configs.Config{
		JWTSecret: jwtSecret,
	})
	return NeonAuthMiddleware(neonClient, userRepo)
}
