package ui

import (
	"time"

	"github.com/Rhymen/go-whatsapp"
	"github.com/Rhymen/go-whatsapp/binary/proto"
	"github.com/rivo/tview"
	log "github.com/sirupsen/logrus"
	"github.com/tech-nico/whatsapp-cli/client"
)

type WhatsappHandler struct {
	c   *client.WhatsappClient
	txt *tview.TextView
	app *tview.Application
}

//HandleError needs to be implemented to be a valid WhatsApp handler
func (wa *WhatsappHandler) HandleError(err error) {

	if e, ok := err.(*whatsapp.ErrConnectionFailed); ok {
		log.Printf("Connection failed, underlying error: %v", e.Err)
		log.Println("Waiting 30sec...")
		<-time.After(30 * time.Second)
		log.Println("Reconnecting...")
		err := wa.c.WaC.Restore()
		if err != nil {
			log.Fatalf("Restore failed: %v", err)
		}
	} else {
		log.Warnf("error occurred: %s\n", err)
	}
}

func (wa *WhatsappHandler) ShouldCallSynchronously() bool {
	return true
}

func (wa *WhatsappHandler) HandleTextMessage(message whatsapp.TextMessage) {
	log.Debug("waHandler.HandleTextMessage")
	log.Tracef("%s", message)

}

func (wa *WhatsappHandler) HandleRawMessage(message *proto.WebMessageInfo) {
	log.Debug("In waHander.HanldeRawMessage")
	log.Tracef("%s", message)
}

func NewWhatsappHandler(app *tview.Application, view *tview.TextView) *WhatsappHandler {
	return &WhatsappHandler{
		txt: view,
		app: app,
	}
}
