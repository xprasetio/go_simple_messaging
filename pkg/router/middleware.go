package router

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/kooroshh/fiber-boostrap/app/repository"
	"github.com/kooroshh/fiber-boostrap/pkg/jwt_token"
	"github.com/kooroshh/fiber-boostrap/pkg/response"
)

func MiddlewareValidateAuth(ctx *fiber.Ctx) error {
	auth := ctx.Get("authorization")
	if auth == "" {
		log.Println("authorization header is missing")
		return response.SendFailureResponse(ctx, fiber.StatusUnauthorized, "authorization header is missing", nil)
	}

	_, err := repository.GetUserSessionByToken(ctx.Context(), auth)
	if err != nil {
		log.Println("Session not found for token:", auth, "Error:", err)
		return response.SendFailureResponse(ctx, fiber.StatusUnauthorized, "Session not found for token", nil)
	}
	claim, err := jwt_token.ValidateToken(ctx.Context(), auth)
	if err != nil {
		fmt.Println(err)
		return response.SendFailureResponse(ctx, fiber.StatusUnauthorized, "Invalid Token di Middleware", nil)
	}
	if time.Now().Unix() > claim.ExpiresAt.Unix() {
		log.Println("token has expired", claim.ExpiresAt)
		return response.SendFailureResponse(ctx, fiber.StatusUnauthorized, "token has expired", nil)
	}
	ctx.Locals("username", claim.Username)
	ctx.Locals("fullname", claim.Fullname)

	return ctx.Next()
}

func MiddlewareRefreshToken(ctx *fiber.Ctx) error {
	auth := ctx.Get("authorization")
	if auth == "" {
		log.Println("Authorization header is missing")
		return response.SendFailureResponse(ctx, fiber.StatusUnauthorized, "Authorization header is missing", nil)
	}
	claim, err := jwt_token.ValidateToken(ctx.Context(), auth)
	if err != nil {
		fmt.Println(err)
		return response.SendFailureResponse(ctx, fiber.StatusUnauthorized, "Invalid Token di Middleware", nil)
	}
	if time.Now().Unix() > claim.ExpiresAt.Unix() {
		log.Println("Token has expired", claim.ExpiresAt)
		return response.SendFailureResponse(ctx, fiber.StatusUnauthorized, "Token has expired", nil)
	}
	ctx.Locals("username", claim.Username)
	ctx.Locals("fullname", claim.Fullname)

	return ctx.Next()
}
