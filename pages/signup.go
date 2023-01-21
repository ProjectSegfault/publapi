package pages

import (
	"github.com/gofiber/fiber/v2"
	"os"

	"github.com/containrrr/shoutrrr"

	log "github.com/sirupsen/logrus"
)

// SignupPage is the signup page handler
func SignupPage(c *fiber.Ctx) error {

	username := c.FormValue("username")
	email := c.FormValue("email")
	ssh := c.FormValue("ssh")
	if username == "" || email == "" || ssh == "" {
		log.Error("username, email and ssh must be filled", username, email, ssh)
		return c.SendStatus(fiber.StatusBadRequest)
	}

	// create user file
	f, err := os.Create("/var/publapi/users/" + username + ".sh")
	if err != nil {
		log.Error("Error creating user file", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	defer f.Close()
	bashscript := "#!/bin/bash \n" +
		"# Path: /var/publapi/users/" + username + ".sh\n" +
		"# This file is generated by publapi. Do not edit this file.\n" +
		"echo \"email of " + username + " is " + email + "\"\n" +
		"pass=\"$(tr -dc A-Za-z0-9 </dev/urandom | head -c 64)\"\n" +
		"useradd -Um -s /bin/bash " + username + "\n" +
		"chmod 711 /home/" + username + "\n" +
		"printf \"%s\\n%s\" \"${pass}\" \"${pass}\" | passwd " + username + "\n" +
		"mkdir /home/" + username + "/.ssh\n" +
		"echo '" + ssh + "' > /home/" + username + "/.ssh/authorized_keys\n" +
		"chmod 700 /home/" + username + "/.ssh\n" +
		"chmod 600 /home/" + username + "/.ssh/authorized_keys\n" +
		"chown -R " + username + ":" + username + " /home/" + username + "/.ssh\n" +
		"echo \"${pass}\" > /home/" + username + "/pass\n" +
		"chmod 600 /home/" + username + "/pass\n" +
		"chown " + username + ":" + username + " /home/" + username + "/pass\n" +
		"sed -i 's/REPLACEME/" + username + "/g' /home/" + username + "/{meta-info.env,Caddyfile}\n" +
		"sed -i 's/EMAIL=/EMAIL=" + email + "/' /home/" + username + "/meta-info.env\n" +
		"loginctl enable-linger " + username + "\n" +
		"setquota -u " + username + " 20G 20G 0 0 /\n" +
		"echo \"" + username + "'s account has been created!\"\n" +
		"rm -rf $0"

	chmoderr := os.Chmod("/var/publapi/users/"+username+".sh", 0700)
	if chmoderr != nil {
		log.Error(err)
	}
	// write to file
	_, err = f.WriteString(bashscript)
	if err != nil {
		log.Error("Error writing to user file", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	log.Info("Registration request for " + username + " has been submitted by the frontend and has been written to /var/publapi/users/" + username + ".sh")
	// send notification to admins
	err = shoutrrr.Send(os.Getenv("PUBLAPI_SHOUTRRRURL"), "New user signup! Please review /var/publapi/users/"+username+".sh to approve or deny the user.")
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
