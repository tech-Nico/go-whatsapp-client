package cliprompt

import (
	"fmt"

	go_prompt "github.com/c-bata/go-prompt"
	log "github.com/sirupsen/logrus"
	wc "github.com/tech-nico/whatsapp-cli/client"
)

type Completer struct {
	Client *wc.WhatsappClient
	chats  []string
}

func NewCompleter() *Completer {
	c, err := wc.NewClient()

	if err != nil {

		log.WithFields(log.Fields{
			"event": "login",
			"error": err,
		}).Error("Error while logging in to Whatsapp")
	}

	fmt.Printf("\nGot a new whatsapp client: %s", c.Session)

	return &Completer{
		Client: &c,
	}
}

func (c *Completer) CompleteCommand(doc go_prompt.Document) []go_prompt.Suggest {
	s := []go_prompt.Suggest{
		{Text: "login", Description: "Login into whatsapp scanning a QR code"},
		{Text: "get", Description: "List all the current chats"},
	}
	return go_prompt.FilterHasPrefix(s, doc.GetWordBeforeCursor(), true)
}
