package config

import (
	"github.com/go-ini/ini"
	"gopkg.in/gomail.v2"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type User struct {
	Usename  string
	Password string
	Email    string
}

type Manager struct {
	Email    string
	Password string
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func CheckResponseForSentence(data []byte, sentence string) bool {
	if len(data) > 0 {
		return strings.Contains(string(data), sentence)
	} else {
		return false
	}
}

func CheckStatusOnWeb(IniPath string) {
	const myurl = "https://dirstudio.laziodisco.it"
	users := ReadUsers(IniPath)
	var message string
	var manager *Manager

	for _, section := range users {
		if section.Name() == "Manager" {
			manager = &Manager{section.Key("email").String(), section.Key("password").String()}
		}
		if section.Name() == "User" {
			user := &User{section.Key("username").String(), section.Key("password").String(), section.Key("email").String()}

			client := &http.Client{}

			formData := url.Values{}
			formData.Add("username", user.Usename)
			formData.Add("password", user.Password)

			responce, err := client.PostForm(myurl, formData)
			CheckError(err)
			defer responce.Body.Close()

			data, _ := ioutil.ReadAll(responce.Body)
			targerSentence := "Elenco domande e stato compilazione Anno Accademico 2023/2024"

			if CheckResponseForSentence(data, targerSentence) {
				message = "The user exist in the database of the laziodisco and applied for the scholarship of 2023/2024."
			} else {
				message = "The user does not exist in the database of the laziodisco or did not apply for the scholarship of 2023/2024."
			}

			SendEmail(manager, user, message)
		}

	}

}

func SendEmail(manager *Manager, user *User, TextMessage string) {
	message := gomail.NewMessage()
	message.SetHeader("From", manager.Email)
	message.SetHeader("To", user.Email)
	message.SetHeader("Subject", "Lazio account status")
	message.SetBody("text/plain", TextMessage)

	d := gomail.NewDialer("smtp.gmail.com", 587, manager.Email, manager.Password)
	if err := d.DialAndSend(message); err != nil {
		log.Fatal(err)
	}
}

func WriteUser(user *User) {
	cfg, err := ini.Load("config.ini")
	CheckError(err)

	section, err := cfg.NewSection("User")
	CheckError(err)

	section.NewKey("email", user.Email)
	section.NewKey("username", user.Usename)
	section.NewKey("password", user.Password)

	cfg.SaveTo("config.ini")

}

func ReadUsers(IniPath string) map[string]*ini.Section {

	cfg, err := ini.Load(IniPath)
	CheckError(err)
	userSections := make(map[string]*ini.Section)

	for _, section := range cfg.Sections() {
		userSections[section.Name()] = section
	}
	return userSections
}
