package middleware

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/izzy-Ti/ZemlyGo/internals/dto"
)

func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("[PANIC RECOVERED] %v\nStack:\n%s\n", err, debug.Stack())
				c.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorAPIResponse("Internal server error"))
			}
		}()
		c.Next()
	}
}
