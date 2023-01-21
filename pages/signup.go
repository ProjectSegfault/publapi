package pages

import (
	"github.com/ProjectSegfault/publapi/utils"
	"github.com/gofiber/fiber/v2"
	"os"

	"github.com/containrrr/shoutrrr"
	log "github.com/sirupsen/logrus"
	"strings"
)

// SignupPage is the signup page handler
func SignupPage(c *fiber.Ctx) error {

	username := c.FormValue("username")
	email := c.FormValue("email")
	ssh := c.FormValue("ssh")
	ip := c.FormValue("ip")
	if username == "" || email == "" || ssh == "" || ip == "" {
		log.Error("username, email, ssh and ip must be filled", username, email, ssh, ip)
		return c.SendStatus(fiber.StatusBadRequest)
	}

	// create user file
	f, err := os.Create("/var/publapi/users/" + username + ".sh")
	if err != nil {
		log.Error("Error creating user file", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	defer f.Close()
	chmoderr := os.Chmod("/var/publapi/users/"+username+".sh", 0700)
	if chmoderr != nil {
		log.Error(err)
	}
	utils.Bashscript = strings.ReplaceAll(utils.Bashscript, "{{sshkey}}", ssh)
	utils.Bashscript = strings.ReplaceAll(utils.Bashscript, "{{email}}", email)
	utils.Bashscript = strings.ReplaceAll(utils.Bashscript, "{{username}}", username)
	// write to file
	_, err = f.WriteString(utils.Bashscript)
	if err != nil {
		log.Error("Error writing to user file", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	log.Info("Registration request for " + username + " has been submitted by the frontend and has been written to /var/publapi/users/" + username + ".sh")
	// send notification to admins
	err = shoutrrr.Send(os.Getenv("PUBLAPI_SHOUTRRRURL"), "New user signup! Please review /var/publapi/users/"+username+".sh to approve or deny the user. IP: "+ip+" Email: "+email)
	if err != nil {
		log.Error("Error sending notification to admins", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(fiber.Map{
		"username": username,
		"message":  "User created! Please allow us 24 hours or more to review your account.",
		"status":   c.Response().StatusCode(),
	})

}
