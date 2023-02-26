package pages

import (
	"github.com/ProjectSegfault/publapi/utils"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
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
	Op        bool   `json:"op"`
	Created   int    `json:"created"`
	Email     string `json:"email"`
	Matrix    string `json:"matrix"`
	Fediverse string `json:"fediverse"`
	Website   string `json:"website"`
	Capsule   string `json:"capsule"`
	Loc       string `json:"loc"`
}

func userdata(username, usersonline, ops string) Userinfo {
	regex := "(^| )" + username + "($| )"
	isonline, err := regexp.MatchString(string(regex), string(usersonline))
	if err != nil {
		log.Error(err)
	}
	isop, operr := regexp.MatchString(string(regex), string(ops))
	if operr != nil {
		log.Error(err)
	}
	cmd := "/usr/bin/stat -c %W /home/" + username
	crd, crerr := exec.Command("bash", "-c", cmd).Output()
	if crerr != nil {
		log.Error(crerr)
	}
	crdstr := string(crd)
	crdstr = strings.TrimSuffix(crdstr, "\n")
	filename := "/home/" + username + "/meta-info.toml"
	_, error := os.Stat(filename)
	if error != nil {
		if os.IsNotExist(error) {
			log.Error(username + " does not have a meta-info.toml")
			var user Userinfo
			user.Name = username
			user.Created, _ = strconv.Atoi(crdstr)
			if isonline == true {
				user.Online = true
			} else {
				user.Online = false
			}
			if isop == true {
				user.Op = true
			} else {
				user.Op = false
			}
			return user
		}
	}
	viper.SetConfigFile(filename)
	viper.ReadInConfig()
	var user Userinfo
	user.Name = username
	user.Created, _ = strconv.Atoi(crdstr)
	user.FullName = viper.GetString("fullname")
	user.Capsule = viper.GetString("gemini")
	user.Website = viper.GetString("website")
	user.Desc = viper.GetString("description")
	user.Email = viper.GetString("email")
	user.Matrix = viper.GetString("matrix")
	user.Fediverse = viper.GetString("fediverse")
	user.Loc = viper.GetString("location")
	if isop == true {
		user.Op = true
	} else {
		user.Op = false
	}
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
	var output int
	if strings.Contains(usersonlinededup, " ") {
		outputa := int(strings.Count(usersonlinededup, " "))
		output = outputa + 1
	} else if usersonlinededup == "" {
		output = 0
	} else {
		output = 1
	}
	// Get OPs
	ops, opserr := exec.Command("bash", "-c", "/usr/bin/members sudo").Output()
	if opserr != nil {
		log.Error(err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	opstr := string(ops)
	// Get all users
	users, err2 := exec.Command("bash", "-c", "/usr/bin/ls /home").Output()
	if err2 != nil {
		log.Error(err2)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	userstr := string(users)
	userstr2 := strings.TrimSuffix(userstr, "\n")
	usersarr := strings.Split(userstr2, "\n")
	// Fill user info
	var userinfostruct []Userinfo
	for i := 0; i < len(usersarr); i++ {
		uname := string(usersarr[i])
		userinfostruct = append(userinfostruct, userdata(uname, strings.TrimSuffix(usersonlinededup, "\n"), strings.TrimSuffix(opstr, "\n")))
	}
	data := Userstruct{
		Status: c.Response().StatusCode(),
		Online: output,
		Users:  userinfostruct,
		Total:  int(strings.Count(userstr, "\n")),
	}
	return c.JSON(data)
}
