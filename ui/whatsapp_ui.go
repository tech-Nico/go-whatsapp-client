package ui

import (
	"fmt"

	"github.com/Rhymen/go-whatsapp"
	"github.com/Rhymen/go-whatsapp/binary/proto"
	log "github.com/sirupsen/logrus"
	"github.com/tech-nico/whatsapp-cli/client"
	"gitlab.com/tslocum/cview"
)

type UI struct {
	Client          client.WhatsappClient
	Pages           *cview.Pages
	App             *cview.Application
	LogView         *cview.TextView
	ChatView        *cview.TextView
	selectedContact whatsapp.Chat
	ContactList     *cview.List
	myJID           string
	cyclePrimitives []*cview.Primitive //This is used to cycle primitives focus when user press Tab
	imagesIDs       []string           //We'll store images IDs here so we can do the highlighting and show them when the user chose to show an image
}

func (thisUI *UI) BuildInfoBar() *cview.TextView {
	// The bottom row has some info on where we are.
	info := cview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(false).
		SetHighlightedFunc(func(added, removed, remaining []string) {
			thisUI.Pages.SwitchToPage(added[0])
		})
	fmt.Fprintf(info, `%d ["%s"][darkcyan]%s[white][""]  `, 1, "chats-page", "Chats")
	fmt.Fprintf(info, `%d ["%s"][darkcyan]%s[white][""]  `, 2, "logs-page", "Logs")

	return info

}

func ShowApp() (*UI, error) {
	thisUI := &UI{}

	thisUI.App = cview.NewApplication()

	thisUI.Pages = cview.NewPages()

	flex, err := thisUI.BuildChatWindow()

	if err != nil {
		return &UI{}, err
	}

	logView := thisUI.BuildLogView()

	thisUI.Pages.AddPage("chats-page", flex, true, true)
	thisUI.Pages.AddPage("logs-page", logView, true, false)

	// The bottom row has some info on where we are.
	info := thisUI.BuildInfoBar()

	// Create the main layout.
	layout := cview.NewFlex().
		SetDirection(cview.FlexRow).
		AddItem(thisUI.Pages, 0, 1, true).
		AddItem(info, 1, 1, false)

	if err := thisUI.App.SetRoot(layout, true).EnableMouse(true).Run(); err != nil {
		return &UI{}, err
	}

	return thisUI, nil
}

func (thisUI *UI) SendMessage(msg interface{}) {
	fmt.Print("\nSend msg to chat window")
	log.Debugf("Sending message to chat window")
	switch msg := msg.(type) {
	case *proto.WebMessageInfo:
		fmt.Print("\nget proto.WebMessageInfo")
		log.Debugf("doGetContacts: received raw message: %s", msg)
	case whatsapp.TextMessage:
		fmt.Printf("\nGot whatsapp.TextMessage: %v", msg)
		thisUI.App.QueueUpdateDraw(func() {
			fmt.Fprintf(thisUI.ChatView, "%s\n%s\n", client.FormatDate(msg.Info.Timestamp), msg.Text)
		})
	default:
		fmt.Printf("\vMessage type unknown: %T", msg)
		return

	}

}
