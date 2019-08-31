package main

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	whatsappclient "github.com/tech-nico/go-whatsapp-client/whatsappclient"
)

func initLogs() {
	log.SetLevel(log.TraceLevel)
	log.SetFormatter(&log.TextFormatter{})
	log.SetReportCaller(true)
}

func main() {
	initLogs()
	c, err := whatsappclient.NewClient()
	if err != nil {
		log.WithFields(log.Fields{
			"event": "login",
			"error": err,
		}).Error("Error while logging in to Whatsapp")
	}

	fmt.Printf("\nGot a new whatsapp client: %s", c.Session)
}
