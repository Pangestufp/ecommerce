package middleware

import (
	"backend/errorhandler"

	"github.com/gin-gonic/gin"
)

func RoleMiddleware(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get("role")
		if !exists {
			errorhandler.ErrorHandler(c, &errorhandler.ForbiddenError{Message: "You cannot access this API"})
			c.Abort()
			return
		}

		if roleVal == role {
			c.Next()
			return
		}

		errorhandler.ErrorHandler(c, &errorhandler.ForbiddenError{Message: "You cannot access this API"})
		c.Abort()
	}
}
