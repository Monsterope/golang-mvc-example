package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type App struct {
	App *fiber.App
}

func NewApp() *App {
	fiberConApp := &App{
		App: fiber.New(),
	}

	fiberConApp.App.Use(cors.New())

	return fiberConApp

}

func (fiberConApp *App) Start(address string) {
	fiberConApp.App.Listen(address)
}
