package cliprompt

import (
	"fmt"

	"github.com/c-bata/go-prompt"
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

//CompleteCommand complete a command
func (c *Completer) CompleteCommand(doc go_prompt.Document) []go_prompt.Suggest {
	if doc.TextBeforeCursor() == "" {
		return []prompt.Suggest{}
	}
	s := []go_prompt.Suggest{
		{Text: "login", Description: "Login into whatsapp scanning a QR code"},
		{Text: "logout", Description: "Logout this session"},
		{Text: "get", Description: "List all the current chats"},
	}
	return go_prompt.FilterHasPrefix(s, doc.GetWordBeforeCursor(), true)
}
