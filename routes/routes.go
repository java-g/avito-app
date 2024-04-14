package routes

import (
	"avito-app/controllers"
	"avito-app/middlewares"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {

	app.Post("/register", controllers.PostRegister)
	app.Post("/login", controllers.PostLogin)

	app.Use(middlewares.IsUser)
	app.Get("/user_banner", controllers.GetUserBanner)

	app.Use(middlewares.IsAdmin)
	app.Get("/banner", controllers.GetBanner)
	app.Post("/banner", controllers.PostBanner)
	app.Patch("/banner/:id", controllers.PatchBanner)
	app.Delete("/banner/:id", controllers.DeleteBanner)
}
