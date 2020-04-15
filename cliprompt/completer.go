package cliprompt

import (
	"strings"

	"github.com/Rhymen/go-whatsapp"
	"github.com/c-bata/go-prompt"
	go_prompt "github.com/c-bata/go-prompt"
	log "github.com/sirupsen/logrus"
	wc "github.com/tech-nico/whatsapp-cli/client"
)

type Completer struct {
	Client *wc.WhatsappClient
	Chats  map[string]whatsapp.Chat
}

func NewCompleter() *Completer {
	c, err := wc.NewClient()

	if err != nil {

		log.WithFields(log.Fields{
			"event": "login",
			"error": err,
		}).Error("Error while logging in to Whatsapp")
	}

	chats, err := c.GetChats()

	if err != nil {
		log.Errorf("Error while retriving chats: %s", err)
		return nil
	}

	return &Completer{
		Client: &c,
		Chats:  chats,
	}
}

func (c *Completer) isCommand(p string) bool {
	return len(strings.Split(p, " ")) < 2
}

func (c *Completer) commandsSuggestion(arg string) []prompt.Suggest {
	s := []go_prompt.Suggest{
		{Text: "login", Description: "Login into whatsapp scanning a QR code"},
		{Text: "logout", Description: "Logout this session"},
		{Text: "get", Description: "List all the current chats"},
	}

	filtered := go_prompt.FilterHasPrefix(s, arg, true)
	if len(filtered) > 0 {
		s = filtered
	}

	return s
}

func (c *Completer) getCommand(p string) string {
	words := strings.Split(p, " ")
	return words[0]
}

func (c *Completer) determineWhatToGet(doc go_prompt.Document) string {
	words := strings.Split(doc.TextBeforeCursor(), " ")
	return words[1]
}

func (c *Completer) handleGet(doc go_prompt.Document) []go_prompt.Suggest {
	cmd := c.determineWhatToGet(doc)

	switch cmd {
	case "history":
		return c.handleHistory(doc)
	default:
		suggestions := []go_prompt.Suggest{
			{Text: "chats", Description: "Retrieve all existing chats"},
			{Text: "history", Description: "Get a chat history"},
			{Text: "contacts", Description: "Get the list of contacts"},
			{Text: "chat", Description: "Start chatting with someone or an existing group"},
		}
		filtered := go_prompt.FilterHasPrefix(suggestions, cmd, true)
		if len(filtered) > 0 {
			suggestions = filtered
		}

		return suggestions
	}
}

func (c *Completer) handleHistory(doc go_prompt.Document) []go_prompt.Suggest {
	suggestions := []go_prompt.Suggest{}
	for _, chat := range c.Chats {
		suggestion := go_prompt.Suggest{
			Text: chat.Name,
		}
		suggestions = append(suggestions, suggestion)
	}
	arg := doc.GetWordBeforeCursor()
	filtered := go_prompt.FilterHasPrefix(suggestions, arg, true)
	if len(filtered) > 0 {
		suggestions = filtered
	}

	return suggestions
}

func (c *Completer) handleChat(doc go_prompt.Document) []go_prompt.Suggest {
	return []go_prompt.Suggest{}
}

//CompleteCommand complete a command
func (c *Completer) CompleteCommand(doc go_prompt.Document) []go_prompt.Suggest {
	p := doc.TextBeforeCursor()

	if c.isCommand(p) {
		return c.commandsSuggestion(doc.GetWordBeforeCursor())
	}

	cmd := c.getCommand(p)
	switch cmd {
	case "get":
		return c.handleGet(doc)
	case "chat":
		return c.handleChat(doc)
	default:
		return []prompt.Suggest{}
	}

}
