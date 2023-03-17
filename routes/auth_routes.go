// Package routes file.
package routes

import (
	"github.com/feserr/pheme-backend/controllers"
	"github.com/gofiber/fiber/v2"
)

// AuthSetup creates the authentication routes.
func AuthSetup(app *fiber.App) {
	app.Post("/api/v1/auth/register", controllers.Register)
	app.Post("/api/v1/auth/login", controllers.Login)
	app.Get("/api/v1/auth/user", controllers.User)
	app.Post("/api/v1/auth/logout", controllers.Logout)
	app.Delete("/api/v1/auth/user", controllers.Delete)
}
