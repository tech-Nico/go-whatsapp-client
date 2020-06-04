package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/Rhymen/go-whatsapp"
	"github.com/Rhymen/go-whatsapp/binary/proto"
	log "github.com/sirupsen/logrus"
	"github.com/tech-nico/whatsapp-cli/client"
	"gitlab.com/tslocum/cview"
)

type WhatsappHandler struct {
	ui *UI
	c  *client.WhatsappClient
}

func (wa *WhatsappHandler) getMessageAuthor(msgInfo whatsapp.MessageInfo) string {
	authorID := "-"
	if msgInfo.FromMe {
		authorID = wa.c.WaC.Info.Wid
	} else {
		if msgInfo.Source.Participant != nil {
			authorID = *msgInfo.Source.Participant
		} else {
			authorID = msgInfo.RemoteJid
		}
	}

	return authorID

}

//Get the message header, i.e. contact name and timestamp
func (wa *WhatsappHandler) buildMessageHeader(messageInfo whatsapp.MessageInfo) string {
	authorID := wa.getMessageAuthor(messageInfo)
	screenName := "-"
	//msgId := message.Info.Id
	if messageInfo.FromMe {
		screenName = "Me"
	} else {
		if messageInfo.Source.PushName != nil {
			screenName = *messageInfo.Source.PushName
		}
		if screenName == "-" {
			if contact, ok := wa.c.Contacts[authorID]; ok {
				screenName = contact.Name
			}
		}
	}

	dateFmt := client.FormatDate(messageInfo.Timestamp)
	return fmt.Sprintf("%s (%s)", dateFmt, screenName)

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
		if !strings.Contains(err.Error(), "error processing data: received invalid data") {
			log.Warnf("error occurred: %s", err)
		}

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

func (wa *WhatsappHandler) PrintAnsiMessage(msgInfo whatsapp.MessageInfo, header, txt string) {
	txt = cview.TranslateANSI(txt)
	wa.PrintMessage(msgInfo, header, txt)
}

//PrintMessage print a message with header + text to the Chat WIndowi
func (wa *WhatsappHandler) PrintMessage(msgInfo whatsapp.MessageInfo, header, txt string) {

	if !msgInfo.FromMe {
		//For now we take for granted the header will never be printed using ANSI formatting code
		header = "[::b]" + header
		txt = "[::b]" + txt
	}

	//This logic is completely wrong since we want to use this function to print others messages but also messages sent by me to selectedContact
	if wa.ui.selectedContact.Jid == msgInfo.RemoteJid {
		wa.ui.ChatView.Write([]byte(fmt.Sprintf("%s\n", header)))
		wa.ui.ChatView.Write([]byte(txt))
		wa.ui.ChatView.Write([]byte("\n\n"))
		wa.ui.ChatView.ScrollToEnd()
		wa.ui.App.Draw()
	}

}

func (wa *WhatsappHandler) HandleImageMessage(msg whatsapp.ImageMessage) {
	log.Trace("Got image message")
	id := *msg.Info.Source.Key.Id
	wa.ui.setImagesIDs(append(wa.ui.getImagesIDs(), id))
	header := wa.buildMessageHeader(msg.Info)

	txt := fmt.Sprintf(`["%s"]Highlight to show the image[""]`, id)
	wa.PrintAnsiMessage(msg.Info, header, txt)
}

func (wa *WhatsappHandler) HandleVideoMessage(msg whatsapp.VideoMessage) {
	log.Infof("Not yet handled: Got video message: %v")
}

func (wa *WhatsappHandler) HandleAudioMessage(msg whatsapp.AudioMessage) {
	log.Infof("Not yet handled: Got audio message: %v", msg)
}

func (wa *WhatsappHandler) HandleDocumentMessage(msg whatsapp.DocumentMessage) {
	log.Infof("Not yet handled: Got document message: %v", msg)
}

func (wa *WhatsappHandler) HandleLiveLocationMessage(msg whatsapp.LiveLocationMessage) {
	log.Infof("Not yet handled: Got live location message: %v", msg)
}

func (wa *WhatsappHandler) HandleLocationMessage(msg whatsapp.LocationMessage) {
	log.Infof("Not yet handled: Got location message: %v", msg)
}

func (wa *WhatsappHandler) HandleStickerMessage(msg whatsapp.StickerMessage) {
	log.Infof("Not yet handled: Got sticker message: %v", msg)
}

func (wa *WhatsappHandler) HandleContactMessage(msg whatsapp.StickerMessage) {
	log.Infof("Not yet handled: Got contact message: %v", msg)
}

func (wa *WhatsappHandler) HandleTextMessage(message whatsapp.TextMessage) {

	messageHeader := wa.buildMessageHeader(message.Info)
	messageStr := message.Text
	log.Debugf("waHandler.HandleTextMessage: from %s", messageHeader)
	log.Tracef("%s", messageStr)
	wa.PrintMessage(message.Info, messageHeader, messageStr)
}

func (wa *WhatsappHandler) HandleRawMessage(message *proto.WebMessageInfo) {
	log.Debugf("In waHander.HanldeRawMessage from %s", message.GetPushName())
	log.Tracef("%s", message)
}

func NewWhatsappHandler(ui *UI) *WhatsappHandler {
	handler := &WhatsappHandler{
		ui: ui,
	}

	handler.c = &ui.Client

	return handler
}
