package main

import (
	"os"
	"os/exec"

	"github.com/ProjectSegfault/publapi/pages"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// publapi is a simple API for Project Segfault's public shared server (pubnix).
func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "welcome to publapi",
			"status":  c.Response().StatusCode(),
		})
	})

	app.Get("/online", func(c *fiber.Ctx) error {

		// set up logger
		logger, _ := zap.NewProduction()
		defer logger.Sync()

		// Get the number of users online
		out, err := exec.Command("users | wc -l").Output()
		if err != nil {
			logger.Error("failed to get number of users online", zap.Error(err))
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.JSON(fiber.Map{
			"users":  out,
			"status": c.Response().StatusCode(),
		})
	})

	app.Post("/signup", pages.SignupPage)

	app.Listen(GetPort())
}

// GetPort returns the port to listen on
func GetPort() string {
	port := os.Getenv("PUBLAPI_PORT")
	if port == "" {
		port = "3000"
	}
	return ":" + port
}
