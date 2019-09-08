package main

import (
	"fmt"

	"github.com/c-bata/go-prompt"
	log "github.com/sirupsen/logrus"
	whatsappclient "github.com/tech-nico/go-whatsapp-client/whatsappclient"
)

func initLogs() {
	log.SetLevel(log.TraceLevel)
	log.SetFormatter(&log.TextFormatter{})
	log.SetReportCaller(true)
}

func completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "users", Description: "Store the username and age"},
		{Text: "articles", Description: "Store the article text posted by user"},
		{Text: "comments", Description: "Store the text commented to articles"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func main() {
	initLogs()

	fmt.Println("Please select table.")

	t := prompt.Input("> ", completer)
	//Check prompt.New as in https://github.com/c-bata/kube-prompt/blob/master/main.go#L33
	fmt.Println("You selected " + t)

	c, err := whatsappclient.NewClient()

	if err != nil {

		log.WithFields(log.Fields{
			"event": "login",
			"error": err,
		}).Error("Error while logging in to Whatsapp")
	}

	fmt.Printf("\nGot a new whatsapp client: %s", c.Session)
}
