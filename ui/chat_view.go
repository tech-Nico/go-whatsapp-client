package ui

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/Rhymen/go-whatsapp"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	log "github.com/sirupsen/logrus"
	"github.com/tech-nico/whatsapp-cli/client"
)

func (thisUI *UI) loadChat(chatFrame *tview.Frame, chatView *tview.TextView, chat whatsapp.Chat) func() {
	return func() {

		thisUI.selectedContact = chat
		currItem := thisUI.ContactList.GetCurrentItem()
		currItemTxt, _ := thisUI.ContactList.GetItemText(currItem)
		thisUI.ContactList.SetItemText(currItem, currItemTxt, "")
		chatFrame.SetTitle(fmt.Sprintf("Chat with %s", chat.Name))
		chatView.Clear()
		chatView.SetText("Loading history...")

		messages, err := thisUI.Client.GetHistory(chat.Jid, 100)
		if err != nil {
			log.Errorf("Error while getting whatsapp history for %s: %s", chat.Jid, err)
		}
		chatView.Clear()
		for idx := range messages {
			currMessage := messages[idx]
			switch currMessage.(type) {
			case whatsapp.TextMessage:
				msg := currMessage.(whatsapp.TextMessage)
				msgId := msg.Info.Id
				authorID := "-"
				screenName := "-"
				if msg.Info.FromMe {
					authorID = thisUI.Client.WaC.Info.Wid
					screenName = "Me"
				} else {
					if msg.Info.Source.Participant != nil {
						authorID = *msg.Info.Source.Participant
					} else {
						authorID = msg.Info.RemoteJid
					}
					if msg.Info.Source.PushName != nil {
						screenName = *msg.Info.Source.PushName
					}
					if screenName == "-" {
						if contact, ok := thisUI.Client.Contacts[authorID]; ok {
							screenName = contact.Name
						}
					}
				}
				dateFmt := client.FormatDate(msg.Info.Timestamp)
				messageHeader := fmt.Sprintf("%s (%s)", dateFmt, screenName)
				message := msg.Text

				if !msg.Info.FromMe {
					messageHeader = "[::b]" + messageHeader
					message = "[::b]" + message
				}
				fmt.Fprintf(thisUI.ChatView, "%s\n", messageHeader)
				fmt.Fprintf(thisUI.ChatView, "%s\n\n", message)

				thisUI.Client.WaC.Read(msg.Info.RemoteJid, msgId)

				//h.messages = append(h.messages, fmt.Sprintf("\n%s (%s): \n%s\n", dateFmt, screenName, message.Text))
			}
			//fmt.Fprintf(chatView, "%s\n", messages[idx])
		}
		chatView.ScrollToEnd()

	}
}

//obsolete - remove
/* func (thisUI *UI) buildContactsList(chatFrame *tview.Frame, chatView *tview.TextView) (*tview.Frame, error) {
	chats := thisUI.Client.Chats

	contactsList := tview.NewList()
	contactsList.ShowSecondaryText(true)

	contactsFrame := tview.NewFrame(contactsList)
	contactsFrame.SetBorder(true)
	contactsFrame.AddText("Chats", true, tview.AlignCenter, tcell.ColorGreen)

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
	thisUI.loadChat(chatFrame, chatView, chats[0])()

	return contactsFrame, nil
}
*/

func (thisUI *UI) alignRight(txt string) string {
	_, _, width, _ := thisUI.ChatView.GetInnerRect()
	maxWidth := width - 5
	messageLength := len(txt)
	paddingSpaces := width - int((math.Min(float64(messageLength), float64(maxWidth))))
	message := strings.Repeat(" ", paddingSpaces) + txt
	return message
}

func (thisUI *UI) buildSendMessageInputFrame(chatView *tview.TextView) *tview.Frame {

	inputField := tview.NewInputField().
		SetLabel("Enter to send: ").
		SetPlaceholder("Hello world")

	inputField.SetDoneFunc(func(key tcell.Key) {
		msgStr := client.FormatDate(uint64(time.Now().Unix()))
		messageHeader := thisUI.alignRight(fmt.Sprintf("%s (Me)", msgStr))
		message := thisUI.alignRight(inputField.GetText())
		fmt.Fprintf(thisUI.ChatView, "%s\n", messageHeader)
		fmt.Fprintf(thisUI.ChatView, "%s\n", message)

		msg := whatsapp.TextMessage{
			Info: whatsapp.MessageInfo{
				RemoteJid: thisUI.selectedContact.Jid,
			},
			Text: inputField.GetText(),
		}
		thisUI.Client.WaC.Send(msg)
		inputField.SetText("")
	})
	inputFrame := tview.NewFrame(inputField)
	inputFrame.SetBorder(true).SetTitle("Enter message")

	return inputFrame
}

func (thisUI *UI) BuildChatWindow() (*tview.Flex, error) {

	chatView := tview.NewTextView()
	chatView.SetBorder(true)
	chatView.SetDynamicColors(true)
	chatView.SetWordWrap(true)
	chatView.SetChangedFunc(func() {
		thisUI.App.Draw()
	})

	//Now we have all we need to create a new instance of the Whatsapp client
	handler := NewWhatsappHandler(thisUI)
	client, err := client.NewClient(handler)

	if err != nil {
		return &tview.Flex{}, fmt.Errorf("Error while initializing WhatsappClient: %s", err)
	}
	thisUI.Client = client

	thisUI.ChatView = chatView //Not sure this is the right thing to do initializing a struct variable in a unrelated method

	chatFrame := tview.NewFrame(chatView)
	chatFrame.SetBorder(true).SetTitle("Select a chat")

	chats := thisUI.Client.Chats

	contactsList := tview.NewList()
	thisUI.ContactList = contactsList
	contactsList.ShowSecondaryText(true)
	contactsList.SetHighlightFullLine(true)

	contactsFrame := tview.NewFrame(contactsList)
	contactsFrame.SetBorder(true)
	contactsFrame.AddText("Chats", true, tview.AlignCenter, tcell.ColorGreen)
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
	thisUI.loadChat(chatFrame, chatView, chats[0])()

	inputFrame := thisUI.buildSendMessageInputFrame(chatView)

	chatFlex := tview.NewFlex().
		AddItem(chatFrame, 0, 4, false).
		AddItem(inputFrame, 0, 1, false).
		SetDirection(tview.FlexRow)

	flex := tview.NewFlex().
		AddItem(contactsFrame, 0, 1, true).
		AddItem(chatFlex, 0, 3, false)

	return flex, nil
}
