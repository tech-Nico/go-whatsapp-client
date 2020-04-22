package client

import (
	"fmt"

	whatsapp "github.com/Rhymen/go-whatsapp"
	log "github.com/sirupsen/logrus"
)

// historyHandler for acquiring chat history
type historyHandler struct {
	c        *WhatsappClient
	messages []string
}

func (h *historyHandler) ShouldCallSynchronously() bool {
	return true
}

// handles and accumulates history's text messages.
// To handle images/documents/videos add corresponding handle functions
func (h *historyHandler) HandleTextMessage(message whatsapp.TextMessage) {
	authorID := "-"
	screenName := "-"
	if message.Info.FromMe {
		authorID = h.c.WaC.Info.Wid
		screenName = "Me"
	} else {
		if message.Info.Source.Participant != nil {
			authorID = *message.Info.Source.Participant
		} else {
			authorID = message.Info.RemoteJid
		}
		if message.Info.Source.PushName != nil {
			screenName = *message.Info.Source.PushName
		}
		if screenName == "-" {
			if contact, ok := h.c.contacts[authorID]; ok {
				screenName = contact.Name
			}
		}
	}
	dateFmt := FormatDate(message.Info.Timestamp)
	h.messages = append(h.messages, fmt.Sprintf("\n%s (%s): \n%s\n", dateFmt, screenName, message.Text))

}

func (h *historyHandler) HandleError(err error) {
	log.Printf("Error occured while retrieving chat history: %s", err)
}

//GetHistory Get the history given a chat JID
func (c *WhatsappClient) GetHistory(jid string, count int) []string {
	// create out history handler

	//Load the list of contacts to look up the scren name in the history handler
	if len(c.contacts) == 0 {
		contacts, err := c.GetContacts()

		if err != nil {
			log.Warningf("Error while retrieving contacts: %s", err)
			contacts = make(map[string]whatsapp.Contact)
		}

		c.contacts = contacts
	}
	handler := &historyHandler{c: c}

	// load chat history and pass messages to the history handler to accumulate
	if count <= 0 {
		count = 100
	}
	c.WaC.LoadChatMessages(jid, count, "", true, false, handler)
	//c.wac.LoadFullChatHistory(jid, 300, time.Millisecond*300, handler)
	return handler.messages
}
