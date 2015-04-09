package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/casimir/xdg-go"
	"github.com/tbruyelle/hipchat-go/hipchat"
)

const (
	datetimeFormat = "2006-01-02T15:04:05.999999999-07:00"
	configFileName = "config"
)

var (
	baseURL             *url.URL
	lastKnownDatetime   time.Time
	xdgApp              = xdg.App{Name: "hipchat-message"}
	config              Config
	errorConfigNotFound = fmt.Errorf("config file not found")
	errorConfigBad      = fmt.Errorf("could not parse config file")
)

func init() {
	var err error
	baseURL, err = url.Parse("https://cup1.corp.badoo.com/v2/")
	if err != nil {
		log.Fatal(err)
	}
}

func createNewConfig() {
	var (
		n   int
		err error
	)

	for n = 0; n != 1 || err != nil; {
		fmt.Print("Enter your security token: ")
		n, err = fmt.Scanf("%s", &config.Token)
	}

	for n = 0; n != 1 || err != nil; {
		fmt.Print("Enter your HipChat base url (e.g. https://cup1.corp.badoo.com/v2/): ")
		n, err = fmt.Scanf("%s", &config.BaseURL)
	}

	if err := config.Check(); err != nil {
		log.Fatal(err)
	}

	if err := config.Save(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("config saved, you can use %s as usual\n", filepath.Base(os.Args[0]))
}

func findUsernameByName(c *hipchat.Client, name string) (string, error) {
	var (
		offset int
		users  []hipchat.User
	)

	for {
		u, _, err := c.User.List(offset, 1000, false, false)
		if err != nil {
			return "", err
		}

		if len(u) == 0 {
			break
		}

		users = append(users, u...)

		offset += len(u)
	}

	for _, u := range users {
		if u.Name == name {
			return u.MentionName, nil
		}
	}

	return "", nil

}

func main() {
	var (
		parameterCode            bool
		parameterCreateNewConfig bool
		messageSent              bool
		message                  string
	)

	flag.BoolVar(&parameterCreateNewConfig, "n", false, "create new config file")
	flag.BoolVar(&parameterCode, "c", false, "send message as a code")
	flag.Parse()

	if parameterCreateNewConfig {
		createNewConfig()
		os.Exit(0)
	}

	if err := config.Load(); err != nil {
		fmt.Println("Error while reading config file:", err)
		fmt.Println("If this is is your first time using this app, please run it with -n parameter to create/recreate config.")
		os.Exit(1)
	}

	if err := config.Check(); err != nil {
		log.Fatal(err)
	}

	if flag.NArg() == 0 {
		fmt.Println("Please enter receiver as a parameter in one of the forms:")
		fmt.Println("  1. Name")
		fmt.Println("  2. @username")
		fmt.Println("  3. +room")
		os.Exit(1)
	}

	if flag.NArg() > 1 {
		fmt.Println("We support only one receiver.")
		os.Exit(1)
	}

	recepient := flag.Arg(0)

	c := hipchat.NewClient(config.Token)
	c.BaseURL = config.GetBaseURL()

	messageBuf, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	if parameterCode {
		message = fmt.Sprintf("/code %s", string(messageBuf))
	} else {
		message = string(messageBuf)
	}

	if strings.HasPrefix(recepient, "+") { // sending to room
		room := recepient[1:]

		nr := hipchat.NotificationRequest{
			Message:       message,
			Notify:        true,
			MessageFormat: "text",
		}

		_, err := c.Room.Notification(room, &nr)
		if err != nil {
			log.Fatal(err)
		}
		messageSent = true
	}

	if !strings.HasPrefix(recepient, "@") {
		username, err := findUsernameByName(c, recepient)
		if err != nil {
			log.Fatal(err)
		}

		if username != "" {
			recepient = fmt.Sprintf("@%s", username)
		}
	}

	if strings.HasPrefix(recepient, "@") { // sending to room
		user := recepient // we should include @ character here

		mr := hipchat.MessageRequest{
			Message:       message,
			Notify:        true,
			MessageFormat: "text",
		}

		resp, err := c.User.Message(user, &mr)
		if err != nil {
			log.Printf("%+v\n", resp)
			log.Printf("%+v\n", resp.Request)
			log.Fatal(err)
		}

		messageSent = true
	}

	if !messageSent {
		fmt.Println("Could not find any recepient.")
	}
}
