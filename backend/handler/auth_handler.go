package handler

import (
	"backend/dto"
	"backend/errorhandler"
	"backend/helper"
	"backend/service"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type authHandler struct {
	service service.AuthService
}

func NewAuthHandler(s service.AuthService) *authHandler {
	return &authHandler{
		service: s,
	}
}

func (h *authHandler) RegisterCustomer(c *gin.Context) {
	var register dto.RegisterRequest

	if err := c.ShouldBindJSON(&register); err != nil {
		errorhandler.ErrorHandler(c, &errorhandler.BadRequestError{Message: err.Error()})
		return
	}

	if err := h.service.Register(&register, helper.Customer()); err != nil {
		errorhandler.ErrorHandler(c, err)
		return
	}

	res := helper.BuildResponse(
		dto.ResponseParam{
			StatusCode: http.StatusCreated,
			Message:    "Register Success",
		})

	c.JSON(http.StatusCreated, res)

}

func (h *authHandler) Login(c *gin.Context) {
	var login dto.LoginRequest

	if err := c.ShouldBindJSON(&login); err != nil {
		errorhandler.ErrorHandler(c, &errorhandler.BadRequestError{Message: err.Error()})
		return
	}

	result, newRefreshToken, err := h.service.Login(c.Request.Context(), &login)

	if err != nil {
		errorhandler.ErrorHandler(c, err)
		return
	}

	// ganti secure=true pas production HTTPS
	c.SetCookie(
		"refresh_token",
		*newRefreshToken,
		7*24*60*60,
		"/",
		"",
		false,
		true,
	)

	res := helper.BuildResponse(
		dto.ResponseParam{
			StatusCode: http.StatusOK,
			Message:    "success",
			Data:       result,
		})

	c.JSON(http.StatusOK, res)
}

func (h *authHandler) Refresh(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")

	if !strings.HasPrefix(tokenString, "Bearer ") {
		errorhandler.ErrorHandler(c, &errorhandler.UnauthorizedError{
			Message: "Unauthorized",
		})
		return
	}

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	userID, _, err := helper.ExtractUserFromExpiredToken(tokenString)
	if err != nil {
		errorhandler.ErrorHandler(c, &errorhandler.UnauthorizedError{
			Message: "Unauthorized",
		})
		return
	}

	// ambil refresh token dari cookie
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		errorhandler.ErrorHandler(c, &errorhandler.UnauthorizedError{
			Message: "Refresh token not found",
		})
		return
	}

	result, newRefreshToken, err := h.service.Refresh(
		c.Request.Context(),
		refreshToken,
		*userID,
	)

	if err != nil {
		errorhandler.ErrorHandler(c, err)
		return
	}

	// ganti secure=true pas production HTTPS
	c.SetCookie(
		"refresh_token",
		*newRefreshToken,
		7*24*60*60,
		"/",
		"",
		false,
		true,
	)

	res := helper.BuildResponse(
		dto.ResponseParam{
			StatusCode: http.StatusOK,
			Message:    "success refresh token",
			Data:       result,
		},
	)

	c.JSON(http.StatusOK, res)
}
