package client

import (
	"fmt"
	"time"

	whatsapp "github.com/Rhymen/go-whatsapp"
	log "github.com/sirupsen/logrus"
)

// historyHandler for acquiring chat history
type historyHandler struct {
	c        *whatsapp.Conn
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
		authorID = h.c.Info.Wid
		screenName = ""
	} else {
		if message.Info.Source.Participant != nil {
			authorID = *message.Info.Source.Participant
		} else {
			authorID = message.Info.RemoteJid
		}
		if message.Info.Source.PushName != nil {
			screenName = *message.Info.Source.PushName
		}
	}

	date := time.Unix(int64(message.Info.Timestamp), 0)
	h.messages = append(h.messages, fmt.Sprintf("%s	%s (%s): %s", date,
		authorID, screenName, message.Text))

}

func (h *historyHandler) HandleError(err error) {
	log.Printf("Error occured while retrieving chat history: %s", err)
}

//GetHistory Get the history given a contact ID
func (c *WhatsappClient) GetHistory(jid string) []string {
	// create out history handler
	handler := &historyHandler{c: &c.wac}

	// load chat history and pass messages to the history handler to accumulate
	c.wac.LoadFullChatHistory(jid, 300, time.Millisecond*300, handler)
	return handler.messages
}
