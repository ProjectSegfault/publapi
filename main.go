package main

import (
	"os"
	"os/exec"

	"github.com/ProjectSegfault/publapi/pages"
	"github.com/gofiber/fiber/v2"

	log "github.com/sirupsen/logrus"

	"runtime"
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
		if runtime.GOOS == "windows" {
			return c.JSON(fiber.Map{
				"message": "/online is not supported on Windows",
				"status":  c.Response().StatusCode(),
			})
		}
		// Get the number of users online
		out, err := exec.Command("bash", "-c", "users | wc -l").Output()
		if err != nil {
			log.Error(err)
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
