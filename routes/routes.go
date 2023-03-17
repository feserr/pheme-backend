package routes

import (
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	AuthSetup(app)
	PhemeSetup(app)
	UserSetup(app)
}
