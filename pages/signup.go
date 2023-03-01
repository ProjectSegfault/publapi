package pages

import (
	"github.com/ProjectSegfault/publapi/utils"
	"github.com/containrrr/shoutrrr"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

// SignupPage is the signup page handler
func SignupPage(c *fiber.Ctx) error {
	SignupIP, SignupIPExists := os.LookupEnv("PUBLAPI_SIGNUP_IP")
	if SignupIPExists == true {
		if c.IP() != SignupIP {
			log.Info("Request made from invalid IP: ", c.IP())
			return c.SendStatus(fiber.StatusForbidden)
		}
	}
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
	Bashscript := strings.ReplaceAll(utils.Bashscript, "{{sshkey}}", ssh)
	Bashscript = strings.ReplaceAll(Bashscript, "{{email}}", email)
	Bashscript = strings.ReplaceAll(Bashscript, "{{username}}", username)
	// write to file
	_, err = f.WriteString(Bashscript)
	if err != nil {
		log.Error("Error writing to user file", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	log.Info("Registration request for " + username + " has been submitted by the frontend and has been written to /var/publapi/users/" + username + ".sh")
	// send notification to user that their reg request was sent
	err = shoutrrr.Send(os.Getenv("PUBLAPI_EMAIL_SHOUTRRRURL")+email, "Hello "+username+",\nYour registration request has been sent.\nIt will take a maximum of 48 hours for the request to be processed.\nThank you for being part of the Project Segfault Pubnix.")
	if err != nil {
		log.Error("Error sending email to user", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	// send notification to admins
	err = shoutrrr.Send(os.Getenv("PUBLAPI_NOTIFY_SHOUTRRRURL")+os.Getenv("PUBLAPI_NOTIFY_SHOUTRRRURL"), "New user signup! Please review /var/publapi/users/"+username+".sh to approve or deny the user. IP: "+ip+" Email: "+email)
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
