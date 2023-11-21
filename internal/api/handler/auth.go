package handler

import (
	"api-service/internal/api/request"
	"api-service/internal/service"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type AuthApi struct {
	authService *service.Auth
}

func NewAuthApi(authService *service.Auth) *AuthApi {
	return &AuthApi{
		authService: authService,
	}
}

func (a *AuthApi) Login(c *fiber.Ctx) error {
	var req request.LoginRequest

	err := c.BodyParser(&req)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	isValid, err := service.VerifySignature(req.Address, req.Signature)
	if err != nil || !isValid {
		return fiber.NewError(http.StatusBadRequest, "Bad signature")
	}

	resp, err := a.authService.Login(c.Context(), req.Address)
	if err != nil {
		return err
	}

	return c.JSON(resp)
}

func (a *AuthApi) RefreshToken(c *fiber.Ctx) error {
	var req request.RefreshRequest

	err := c.BodyParser(&req)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	resp, err := a.authService.RefreshToken(c, req.Token)
	if err != nil {
		return err
	}

	return c.JSON(resp)
}
