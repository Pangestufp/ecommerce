package middleware

import (
	"backend/errorhandler"

	"github.com/gin-gonic/gin"
)

func RoleMiddleware(roles []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get("role")
		if !exists {
			errorhandler.ErrorHandler(c, &errorhandler.ForbiddenError{Message: "You cannot access this API"})
			c.Abort()
			return
		}

		role, ok := roleVal.(string)
		if !ok {
			errorhandler.ErrorHandler(c, &errorhandler.ForbiddenError{Message: "You cannot access this API"})
			c.Abort()
			return
		}

		for _, r := range roles {
			if r == role {
				c.Next()
				return
			}
		}

		errorhandler.ErrorHandler(c, &errorhandler.ForbiddenError{Message: "You cannot access this API"})
		c.Abort()
	}
}
