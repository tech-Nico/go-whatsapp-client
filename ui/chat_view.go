package ui

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/Rhymen/go-whatsapp"
	"github.com/gdamore/tcell"
	log "github.com/sirupsen/logrus"
	"github.com/tech-nico/whatsapp-cli/client"
	"gitlab.com/tslocum/cview"
)

func (thisUI *UI) loadChat(chatFrame *cview.Frame, chatView *cview.TextView, chat whatsapp.Chat) func() {
	return func() {

		thisUI.selectedContact = chat
		thisUI.imagesIDs = []string{}

		currItem := thisUI.ContactList.GetCurrentItem()
		currItemTxt, _ := thisUI.ContactList.GetItemText(currItem)
		thisUI.ContactList.SetItemText(currItem, currItemTxt, "")
		chatFrame.SetTitle(fmt.Sprintf("Chat with %s", chat.Name))
		chatView.Clear()
		chatView.Write([]byte("Loading history..."))
		thisUI.App.ForceDraw()

		messages, err := thisUI.Client.GetHistory(chat.Jid, 100)

		if err != nil {
			log.Errorf("Error while getting whatsapp history for %s: %s", chat.Jid, err)
		}
		chatView.Clear()
		wa := NewWhatsappHandler(thisUI)
		wa.SetClient(&thisUI.Client)
		for idx := range messages {

			currMessage := messages[idx]
			switch currMessage.(type) {
			case whatsapp.TextMessage:
				msg := currMessage.(whatsapp.TextMessage)

				wa.HandleTextMessage(msg)
				msgID := msg.Info.Id

				thisUI.Client.WaC.Read(msg.Info.RemoteJid, msgID) //Mark message as read

			case whatsapp.ImageMessage:
				msg := currMessage.(whatsapp.ImageMessage)
				wa.HandleImageMessage(msg)
				msgID := msg.Info.Id
				thisUI.Client.WaC.Read(msg.Info.RemoteJid, msgID) //Mark message as read
			default:
				chatView.Write([]byte(fmt.Sprintf("Message type (%T) not yet handled:", currMessage)))
			}
		}

		chatView.ScrollToEnd()

	}
}

func (thisUI *UI) alignRight(txt string) string {
	_, _, width, _ := thisUI.ChatView.GetInnerRect()
	maxWidth := width - 5
	messageLength := len(txt)
	paddingSpaces := width - int((math.Min(float64(messageLength), float64(maxWidth))))
	message := strings.Repeat(" ", paddingSpaces) + txt
	return message
}

func (thisUI *UI) buildSendMessageInputFrame(chatView *cview.TextView) *cview.Frame {

	inputField := cview.NewInputField().
		SetLabel("Enter to send: ").
		SetPlaceholder("Hello world")

	inputField.SetDoneFunc(func(key tcell.Key) {
		wa := NewWhatsappHandler(thisUI)
		msg := whatsapp.TextMessage{

			Info: whatsapp.MessageInfo{
				RemoteJid: thisUI.selectedContact.Jid,
				FromMe:    true,
				Timestamp: uint64(time.Now().Unix()),
			},
			Text: inputField.GetText(),
		}

		thisUI.Client.WaC.Send(msg)
		header := wa.buildMessageHeader(msg.Info)
		wa.PrintMessage(msg.Info, header, inputField.GetText())
		inputField.SetText("")
	})
	inputFrame := cview.NewFrame(inputField)
	inputFrame.SetBorder(true).SetTitle("Enter message")

	return inputFrame
}

//this function should select based on the timestamp found on the region
func (thisUI *UI) selectMessage(chatView *cview.TextView) func(tcell.Key) {
	return func(key tcell.Key) {
		currentSelection := chatView.GetHighlights()

		if key == tcell.KeyEnter {
			if len(currentSelection) > 0 {
				chatView.Highlight() //Here I should show the image
			} else {
				chatView.Highlight(thisUI.imagesIDs[0]).ScrollToHighlight() //This happens when I press enter the first time
			}
		} else if len(currentSelection) > 0 {
			index := 0
			for idx, id := range thisUI.imagesIDs {
				if id == currentSelection[0] {
					index = idx
					break
				}
			}

			if key == tcell.KeyTab {
				index = (index + 1) % len(thisUI.imagesIDs)
			} else if key == tcell.KeyBacktab {
				index = (index - 1 + len(thisUI.imagesIDs)) % len(thisUI.imagesIDs)
			} else {
				return
			}
			chatView.Highlight(thisUI.imagesIDs[index]).ScrollToHighlight()
		}
	}
}

func (thisUI *UI) BuildChatWindow() (*cview.Flex, error) {

	chatView := cview.NewTextView()
	chatView.SetBorder(true)

	chatView.SetDynamicColors(true).
		SetWordWrap(true).
		SetRegions(true)

	chatView.SetDoneFunc(thisUI.selectMessage(chatView))

	//Now we have all we need to create a new instance of the Whatsapp client
	handler := NewWhatsappHandler(thisUI)
	client, err := client.NewClient(handler)

	if err != nil {
		return &cview.Flex{}, fmt.Errorf("Error while initializing WhatsappClient: %s", err)
	}
	thisUI.Client = client
	thisUI.myJID = client.WaC.Info.Wid
	thisUI.ChatView = chatView //Not sure this is the right thing to do initializing a struct variable in a unrelated method

	chatFrame := cview.NewFrame(chatView)
	chatFrame.SetBorder(true).SetTitle("Select a chat")

	chats := thisUI.Client.Chats

	contactsList := cview.NewList()
	thisUI.ContactList = contactsList
	contactsList.ShowSecondaryText(true)
	contactsList.SetHighlightFullLine(true)

	contactsFrame := cview.NewFrame(contactsList)
	contactsFrame.SetBorder(true)
	contactsFrame.AddText("Chats", true, cview.AlignCenter, tcell.ColorGreen)
	for _, v := range chats {
		unread := ""
		if i, err := strconv.Atoi(v.Unread); err == nil && i > 0 {
			unread = fmt.Sprintf("%d unread messages", i)
		}

		if v.Name != "" {
			contactsList.AddItem(v.Name, unread, 0, thisUI.loadChat(chatFrame, chatView, v))
		} else {
			contactsList.AddItem(v.Jid, unread, 0, thisUI.loadChat(chatFrame, chatView, v))
		}
	}
	//Select the first item in the list
	contactsList.SetCurrentItem(0)
	thisUI.selectedContact = chats[0]
	thisUI.loadChat(chatFrame, chatView, chats[0])()

	inputFrame := thisUI.buildSendMessageInputFrame(chatView)

	chatFlex := cview.NewFlex().
		AddItem(chatFrame, 0, 4, false).
		AddItem(inputFrame, 0, 1, false).
		SetDirection(cview.FlexRow)

	flex := cview.NewFlex().
		AddItem(contactsFrame, 0, 1, true).
		AddItem(chatFlex, 0, 3, false)

	return flex, nil
}
