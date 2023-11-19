package api

import (
	"api-service/internal/api/handler"
	"api-service/internal/config"
	"context"
	"errors"
	"github.com/goccy/go-json"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error,omitempty"`
	Code  int    `json:"code"`
}

func NewFiber(ctx context.Context, jwtCfg *config.Jwt, authHandler *handler.AuthApi) *fiber.App {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		JSONDecoder:           json.Unmarshal,
		JSONEncoder:           json.Marshal,
		ErrorHandler:          errHandler,
	})

	app.Use(func(c *fiber.Ctx) error {
		c.SetUserContext(ctx)
		return c.Next()
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestCompression,
	}))

	api := app.Group("api/v1")

	authGroup := api.Group("/auth")
	authGroup.Post("/login", authHandler.Login)

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			Key: []byte(jwtCfg.SecretKey),
		},
	}))

	return app
}

func errHandler(ctx *fiber.Ctx, err error) error {
	code := http.StatusInternalServerError
	message := "Internal Server Error"

	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
		message = e.Message
	}

	return ctx.Status(code).JSON(NewHttpError(code, message))
}

func NewHttpError(code int, error string) ErrorResponse {
	return ErrorResponse{
		Error: error,
		Code:  code,
	}
}
