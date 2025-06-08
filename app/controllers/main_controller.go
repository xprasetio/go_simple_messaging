package controllers

import "github.com/gofiber/fiber/v2"

func Render(c *fiber.Ctx) error {
	return c.Render("index", nil)
}
