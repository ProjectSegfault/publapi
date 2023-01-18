package main

import (
	"github.com/ProjectSegfault/publapi/pages"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type Userstruct struct {
	Status int `json:"status"`
	Online int `json:"online"`
	Total  int `json:"total"`
	Users  []Userinfo `json:"users"`
}
type Userinfo struct {
	Name     string `json:"name"`
	FullName string `json:"fullName"`
	Loc      string `json:"loc"`
	Email    string `json:"email"`
	Desc     string `json:"desc"`
	Website  string `json:"website"`
	Capsule  string `json:"capsule"`
	Online   bool `json:"online"` 
}

func Dedup(input string) string {
	unique := []string{}

	words := strings.Split(input, " ")
	for _, word := range words {
		// If we alredy have this word, skip.
		if contains(unique, word) {
			continue
		}

		unique = append(unique, word)
	}

	return strings.Join(unique, " ")
}

func contains(strs []string, str string) bool {
	for _, s := range strs {
		if s == str {
			return true
		}
	}
	return false
}

func userdata(username, usersonline string) Userinfo {
	filename := "/home/" + username + "/meta-info.env"
	_, error := os.Stat(filename)
	if error != nil {
		if os.IsNotExist(error) {
			log.Error(username + " does not have a meta-info.env")
			var user Userinfo
			if strings.Contains(usersonline, " "+username) == true {
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
	user.Name = viper.GetString("USERNAME")
	user.FullName = viper.GetString("FULL_NAME")
	user.Capsule = viper.GetString("GEMINI_CAPSULE")
	user.Website = viper.GetString("WEBSITE")
	user.Desc = viper.GetString("DESCRIPTION")
	user.Email = viper.GetString("EMAIL")
	user.Loc = viper.GetString("LOCATION")
	if strings.Contains(usersonline, " "+username) == true {
		user.Online = true
	} else {
		user.Online = false
	}
	return user
}

// publapi is a simple API for Project Segfault's public shared server (pubnix).
func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "welcome to publapi",
			"status":  c.Response().StatusCode(),
		})
	})

	app.Get("/users", func(c *fiber.Ctx) error {
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
		usersonlinededup := Dedup(usersonlinestr)
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
		//var userinfoarr []interface{}
		var userinfostruct []Userinfo
		for i := 0; i < len(usersarr); i++ {
			uname := string(usersarr[i])
			userinfostruct = append(userinfostruct, userdata(uname, usersonlinededup))
		}
		log.Info(userinfostruct)
		data := Userstruct{
			Status: c.Response().StatusCode(),
			Online: output,
			Users:  userinfostruct,
			Total:  int(strings.Count(userstr, "\n")),
		}
		return c.JSON(data)
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
