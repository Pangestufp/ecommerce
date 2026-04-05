package middleware

import (
	"backend/errorhandler"
	"backend/helper"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if strings.HasPrefix(tokenString, "Bearer ") {
			tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		} else {
			errorhandler.ErrorHandler(c, &errorhandler.UnauthorizedError{Message: "Unauthorized"})
			c.Abort()
			return
		}

		if tokenString == "" {
			errorhandler.ErrorHandler(c, &errorhandler.UnauthorizedError{Message: "Unauthorized"})
			c.Abort()
			return
		}

		userID, role, err := helper.ValidateToken(tokenString)
		if err != nil {
			errorhandler.ErrorHandler(c, &errorhandler.UnauthorizedError{Message: err.Error()})
			c.Abort()
			return
		}

		c.Set("role", role)
		c.Set("userID", *userID)
		c.Next()
	}
}
