package whatsappclient

import (
	"fmt"

	"github.com/c-bata/go-prompt"
	log "github.com/sirupsen/logrus"
)

type Completer struct {
	client *WhatsappClient
}

func NewCompleter() *Completer {
	c, err := NewClient()

	if err != nil {

		log.WithFields(log.Fields{
			"event": "login",
			"error": err,
		}).Error("Error while logging in to Whatsapp")
	}

	fmt.Printf("\nGot a new whatsapp client: %s", c.Session)

	return &Completer{
		client: &c,
	}
}

func (c *Completer) Complete(doc prompt.Document) []prompt.Suggest {
	return make([]prompt.Suggest, 0)
}
