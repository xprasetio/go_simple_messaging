package controllers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/kooroshh/fiber-boostrap/app/repository"
	"github.com/kooroshh/fiber-boostrap/pkg/response"
)

func GetMessageHistory(ctx *fiber.Ctx) error {
	resp, err := repository.GetMessageHistory(ctx.Context())
	if err != nil {
		fmt.Println(err)
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, "Failed to retrieve message history", nil)
	}
	return response.SendSuccessResponse(ctx, resp)
}
