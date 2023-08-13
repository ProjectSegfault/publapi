package pages

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/ProjectSegfault/publapi/utils"
	"github.com/containrrr/shoutrrr"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
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
	captchaResp := c.FormValue("h-captcha-response")

	if captchaResp == "" {
		log.Error("Nice try, but the registration won't work unless you answer the captcha.")
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if username == "" || email == "" || ssh == "" || ip == "" {
		log.Error("username, email, ssh and ip must be filled", username, email, ssh, ip)
		return c.SendStatus(fiber.StatusBadRequest)
	}
	raid, ok := os.LookupEnv("PUBLAPI_RAID_MODE")
	if !ok || raid == "1" {
		log.Warn(
			"PUBLAPI_RAID_MODE is on, accepting every request as OK and not doing anything...\n User info: ",
			username,
			" ",
			email,
			" ",
			ip,
			" ",
		)
		return c.SendStatus(fiber.StatusOK)
	}

	// Check the captcha validation.

	captchaSecret, ok := os.LookupEnv("PUBLAPI_CAPTCHA_SECRET")

	params := url.Values{}
	params.Add("response", captchaResp)
	params.Add("secret", captchaSecret)
	body := strings.NewReader(params.Encode())

	req, err := http.NewRequest("POST", "https://hcaptcha.com/siteverify", body)
	if err != nil {
		// handle err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error("Something went wrong fetching the HCatpcha API: ", err)
	}
	defer resp.Body.Close()

	bod, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("Error reading captcha response body", err)
	}
	sb := string(bod)
	log.Info("Captcha response: ", sb)

	type CaptchaResponse struct {
		Success bool `json:"success"`
	}
	var captchaResponse CaptchaResponse
	err = json.Unmarshal([]byte(sb), &captchaResponse)
	if err != nil {
		log.Error("Error unmarshalling captcha response", err)
	}

	if captchaResponse.Success == false {
		log.Error("Captcha validation failed")
		return c.JSON(fiber.Map{
			"username": username,
			"message":  "Sorry! But the captcha validation failed. Please try again.",
			"status":   c.Response().StatusCode(),
		})
	} else {
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

		log.Info(
			"Registration request for " + username + " has been submitted by the frontend and has been written to /var/publapi/users/" + username + ".sh",
		)

		// send notification to user that their reg request was sent

		err = shoutrrr.Send(os.Getenv("PUBLAPI_EMAIL_SHOUTRRRURL")+email, "Hello "+username+",\nYour registration request has been sent.\nIt will take a maximum of 48 hours for the request to be processed.\nThank you for being part of the Project Segfault Pubnix.")
		if err != nil {
			log.Error("Error sending email to user", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		// send notification to admins

		//shoutrrrUrl := os.Getenv("PUBLAPI_NOTIFY_SHOUTRRRURL") + os.Getenv("PUBLAPI_NOTIFY_ROOMS")
		shoutrrrUrl := os.Getenv("PUBLAPI_EMAIL_SHOUTRRRURL")+"contact@projectsegfau.lt"
		err = shoutrrr.Send(
			shoutrrrUrl,
			"New user signup! Please review /var/publapi/users/"+username+".sh to approve or deny the user. IP: "+ip+" Email: "+email,
		)
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
}
