package routes

import (
	"main/app/handlers/traderjoe"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

func JwtContextKey() string {
	return "jwt"
}

func LaunchpadRoute(c *fiber.App) {
	route := c.Group("/")
	route.Post(":launchpad/login", traderjoe.Token)
	c.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte("secret")},
	}))
	route.Post(":launchpad/signature", traderjoe.Signature)
}
