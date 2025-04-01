package routes

import (
	"monsterloveshop/config"
	"monsterloveshop/controllers"
	"monsterloveshop/middleware"

	"github.com/gofiber/fiber/v2"
)

type Route struct {
	App         *fiber.App
	Controllers *controllers.Controller
	Middlewares *middleware.RedisAuthMiddleware
}

type ResponseExJson struct {
	Status   int
	Messages string
}

func NewRoute(app *fiber.App, ctr *controllers.Controller, middleware *middleware.RedisAuthMiddleware) *Route {
	return &Route{
		App:         app,
		Controllers: ctr,
		Middlewares: middleware,
	}
}

func ExampleRouteJson(a *fiber.App) {
	a.Get("/json", func(c *fiber.Ctx) error {
		appName := config.GetEnv("app.name")
		respData := ResponseExJson{
			Status:   200,
			Messages: appName,
		}

		return c.JSON(respData)
	})
}

func ExampleRoute(a *fiber.App) {
	a.Get("/", func(c *fiber.Ctx) error {
		return c.SendString(config.GetEnv("app.name"))
	})
}
