package pages

import (
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"

	"github.com/ProjectSegfault/publapi/utils"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Userstruct struct {
	Status int        `json:"status"`
	Online int        `json:"online"`
	Total  int        `json:"total"`
	Users  []Userinfo `json:"users"`
}

type Userinfo struct {
	Name      string `json:"name"`
	FullName  string `json:"fullName"`
	Desc      string `json:"desc"`
	Online    bool   `json:"online"`
	Email     string `json:"email"`
	Matrix    string `json:"matrix"`
	Fediverse string `json:"fediverse"`
	Website   string `json:"website"`
	Capsule   string `json:"capsule"`
	Loc       string `json:"loc"`
}

func userdata(username, usersonline string) Userinfo {
	filename := "/home/" + username + "/meta-info.env"
	_, error := os.Stat(filename)
	regex := "(^| )" + username + "($| )"
	isonline, err := regexp.MatchString(string(regex), string(usersonline))
	if err != nil {
		log.Error(err)
	}
	if error != nil {
		if os.IsNotExist(error) {
			log.Error(username + " does not have a meta-info.env")
			var user Userinfo
			user.Name = username
			if isonline == true {
				user.Online = true
			} else {
				user.Online = false
			}
			return user
		}
	}
	viper.SetConfigFile(filename)
	viper.ReadInConfig()
	var user Userinfo
	user.Name = username
	user.FullName = viper.GetString("FULL_NAME")
	user.Capsule = viper.GetString("GEMINI_CAPSULE")
	user.Website = viper.GetString("WEBSITE")
	user.Desc = viper.GetString("DESCRIPTION")
	user.Email = viper.GetString("EMAIL")
	user.Matrix = viper.GetString("MATRIX")
	user.Fediverse = viper.GetString("FEDIVERSE")
	user.Loc = viper.GetString("LOCATION")
	if isonline == true {
		user.Online = true
	} else {
		user.Online = false
	}
	return user
}

func UsersPage(c *fiber.Ctx) error {
	if runtime.GOOS == "windows" {
		return c.JSON(fiber.Map{
			"message": "/users is not supported on Windows",
			"status":  c.Response().StatusCode(),
		})
	}
	// Get the number of users online
	usersonline, err := exec.Command("bash", "-c", "/usr/bin/users").Output()
	if err != nil {
		log.Error(err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	usersonlinestr := string(usersonline)
	usersonlinededup := utils.Dedup(usersonlinestr)
	outputa := int(strings.Count(usersonlinededup, " "))
	var output int
	output = outputa + 1
	users, err2 := exec.Command("bash", "-c", "/usr/bin/ls /home").Output()
	if err2 != nil {
		log.Error(err2)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	userstr := string(users)
	userstr2 := strings.TrimSuffix(userstr, "\n")
	usersarr := strings.Split(userstr2, "\n")
	var userinfostruct []Userinfo
	for i := 0; i < len(usersarr); i++ {
		uname := string(usersarr[i])
		userinfostruct = append(userinfostruct, userdata(uname, strings.TrimSuffix(usersonlinededup, "\n")))
	}
	data := Userstruct{
		Status: c.Response().StatusCode(),
		Online: output,
		Users:  userinfostruct,
		Total:  int(strings.Count(userstr, "\n")),
	}
	return c.JSON(data)
}
