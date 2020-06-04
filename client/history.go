package client

import (
	"fmt"
	"strings"

	whatsapp "github.com/Rhymen/go-whatsapp"
	log "github.com/sirupsen/logrus"
)

// historyHandler for acquiring chat history
type historyHandler struct {
	c        *WhatsappClient
	messages []interface{}
}

func (h *historyHandler) ShouldCallSynchronously() bool {
	return true
}

// handles and accumulates history's text messages.
// To handle images/documents/videos add corresponding handle functions
func (h *historyHandler) HandleTextMessage(message whatsapp.TextMessage) {
	log.Debug("In historyHandler.HandleTextMessage")
	/*	authorID := "-"
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
				if contact, ok := h.c.Contacts[authorID]; ok {
					screenName = contact.Name
				}
			}
		}
		//dateFmt := FormatDate(message.Info.Timestamp)
		//h.messages = append(h.messages, fmt.Sprintf("\n%s (%s): \n%s\n", dateFmt, screenName, message.Text))
	*/
	h.messages = append(h.messages, message)

}

//HandleImageMessage Add an image message to the list of messages
func (h *historyHandler) HandleImageMessage(message whatsapp.ImageMessage) {
	log.Debug("In historyHandler.HandleImageMessage")

	h.messages = append(h.messages, message)

}

func (h *historyHandler) HandleError(err error) {
	log.Printf("Error occured while retrieving chat history: %s", err)
}

//GetHistory Get the history given a chat JID
func (wc *WhatsappClient) GetHistory(jid string, count int) ([]interface{}, error) {
	// create out history handler

	handler := &historyHandler{c: wc}

	// load chat history and pass messages to the history handler to accumulate
	if count <= 0 {
		count = 100
	}
	err := wc.WaC.LoadChatMessages(jid, count, "", true, false, handler)
	if err != nil {
		return nil, fmt.Errorf("Error retriving %s chat history: %s", jid, err)
	}

	//wc.wac.LoadFullChatHistory(jid, 300, time.Millisecond*300, handler)
	return handler.messages, nil
}

//GetMessage Get a specific message
func (wc *WhatsappClient) GetMessage(jid, msgID string) (interface{}, error) {
	// create out history handler

	handler := &historyHandler{c: wc}
	test, err2 := wc.WaC.LoadMessagesAfter(jid, msgID, 1)
	if err2 != nil {
		log.Warnf("Error while testing: %s", err2)
	} else {
		log.Infof("Success! Got something: %s", test)
	}

	err := wc.WaC.LoadChatMessages(jid, 2, msgID, true, false, handler)
	if err != nil {
		return nil, fmt.Errorf("error retriving message %s for user with jid %s : %s", msgID, jid, err)
	}

	return handler.messages[0], nil
}

func (wc *WhatsappClient) findChat(jid string) (whatsapp.Chat, error) {
	foundChat := whatsapp.Chat{}

	//Search all the chats with the given name. Since you might have two chats with the same name, ask which one to use
	//in case there are two or more with the same name
	for idx := range wc.Chats {
		currJid := wc.Chats[idx].Jid
		if strings.EqualFold(strings.TrimSpace(currJid), strings.TrimSpace(jid)) {
			foundChat = wc.Chats[idx]
		}
	}

	if (foundChat == whatsapp.Chat{}) {
		return foundChat, fmt.Errorf("pushHistoryToCh: no chat with jid %s could be found", jid)
	}

	return foundChat, nil

}
