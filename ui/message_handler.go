package ui

import (
	"fmt"
	"time"

	"github.com/Rhymen/go-whatsapp"
	"github.com/Rhymen/go-whatsapp/binary/proto"
	log "github.com/sirupsen/logrus"
	"github.com/tech-nico/whatsapp-cli/client"
)

type WhatsappHandler struct {
	ui *UI
	c  *client.WhatsappClient
}

//HandleError needs to be implemented to be a valid WhatsApp handler
func (wa *WhatsappHandler) HandleError(err error) {
	wac := wa.c.WaC
	if e, ok := err.(*whatsapp.ErrConnectionFailed); ok {
		log.Printf("Connection failed, underlying error: %v", e.Err)
		log.Println("Waiting 30sec...")
		<-time.After(30 * time.Second)
		log.Println("Reconnecting...")
		err := wac.Restore()
		if err != nil {
			log.Fatalf("Restore failed: %v", err)
		}
	} else {
		log.Warnf("error occurred: %s\n", err)
	}
}

func (wa *WhatsappHandler) SetClient(c *client.WhatsappClient) {
	wa.c = c
}

func (wa *WhatsappHandler) ShouldCallSynchronously() bool {
	return true
}

func (wa *WhatsappHandler) HandleJsonMessage(msg string) {
	log.Infof("Got json message: %s", msg)
}

func (wa *WhatsappHandler) HandleTextMessage(message whatsapp.TextMessage) {

	authorID := "-"
	screenName := "-"
	//msgId := message.Info.Id
	if message.Info.FromMe {
		authorID = wa.c.WaC.Info.Wid
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
			if contact, ok := wa.c.Contacts[authorID]; ok {
				screenName = contact.Name
			}
		}
	}
	dateFmt := client.FormatDate(message.Info.Timestamp)
	messageHeader := fmt.Sprintf("%s (%s)", dateFmt, screenName)
	messageStr := message.Text
	log.Debugf("waHandler.HandleTextMessage: from %s", messageHeader)
	log.Tracef("%s", messageStr)
	if wa.ui.selectedContact.Jid == authorID {
		fmt.Fprintf(wa.ui.ChatView, "%s\n", messageHeader)
		fmt.Fprintf(wa.ui.ChatView, "%s\n", messageStr)
	}

	//	wa.c.WaC.Read(message.Info.RemoteJid, msgId)
}

func (wa *WhatsappHandler) HandleRawMessage(message *proto.WebMessageInfo) {
	log.Debugf("In waHander.HanldeRawMessage from %s", message.GetPushName())
	log.Tracef("%s", message)
}

func NewWhatsappHandler(ui *UI) *WhatsappHandler {
	return &WhatsappHandler{
		ui: ui,
	}
}
