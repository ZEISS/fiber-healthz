package main

import (
	healthz "github.com/zeiss/fiber-healthz"

	"github.com/gofiber/fiber/v3"
)

func main() {
	app := fiber.New()

	app.Get(healthz.Readiness, healthz.New())
	app.Get(healthz.Liveness, healthz.New())
	app.Get(healthz.Healthz, healthz.NewProbe(healthz.Failure))
	app.Get(healthz.Startup, healthz.New())

	err := app.Listen(":3000")
	if err != nil {
		panic(err)
	}
}
