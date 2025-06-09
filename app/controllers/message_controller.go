package controllers

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/kooroshh/fiber-boostrap/app/repository"
	"github.com/kooroshh/fiber-boostrap/pkg/response"
	"go.elastic.co/apm"
)

func GetMessageHistory(ctx *fiber.Ctx) error {
	span, spanCtx := apm.StartSpan(ctx.UserContext(), "GetMessageHistory", "controller")
	defer span.End()
	resp, err := repository.GetMessageHistory(spanCtx)
	if err != nil {
		log.Println(err)
		return response.SendFailureResponse(ctx, fiber.StatusInternalServerError, "Failed to retrieve message history", nil)
	}
	return response.SendSuccessResponse(ctx, resp)
}
