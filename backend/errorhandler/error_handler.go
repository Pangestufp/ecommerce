package errorhandler

import (
	"backend/dto"
	"backend/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandler(c *gin.Context, err error) {
	statusCode := http.StatusInternalServerError
	switch err.(type) {
	case *NotFoundError:
		statusCode = http.StatusNotFound
	case *BadRequestError:
		statusCode = http.StatusBadRequest
	case *InternalServerError:
		statusCode = http.StatusInternalServerError
	case *UnauthorizedError:
		statusCode = http.StatusUnauthorized
	case *ForbiddenError:
		statusCode = http.StatusForbidden
	}

	response := helper.BuildResponse(dto.ResponseParam{
		StatusCode: statusCode,
		Message:    err.Error(),
	})

	c.JSON(statusCode, response)
}
