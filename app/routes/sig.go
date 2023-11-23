package routes

import (
	"main/app/handlers/traderjoe"

	"github.com/gofiber/fiber/v2"
)

func SigRoute(c *fiber.App) {
	route := c.Group("/launchpad")

	route.Post(":signature", traderjoe.Signature)
}
