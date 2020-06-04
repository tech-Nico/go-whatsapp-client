package ui

import (
	"bytes"
	"fmt"
	"image/color"
	"math"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/Rhymen/go-whatsapp"
	"github.com/eliukblau/pixterm/pkg/ansimage"
	"github.com/gdamore/tcell"
	log "github.com/sirupsen/logrus"
	"github.com/tech-nico/whatsapp-cli/client"
	"gitlab.com/tslocum/cview"
)

func (thisUI *UI) setSelectedContactMessages(messages []interface{}) {
	thisUI.RWMutex.Lock()
	defer thisUI.RWMutex.Unlock()
	thisUI.SelectedContactMessages = messages
}

func (thisUI *UI) getSelectedContactMessages() []interface{} {
	thisUI.RWMutex.RLock()
	defer thisUI.RWMutex.RUnlock()
	return thisUI.SelectedContactMessages
}

func (thisUI *UI) setSelectedContact(chat whatsapp.Chat) {
	thisUI.RWMutex.Lock()
	defer thisUI.RWMutex.Unlock()
	thisUI.selectedContact = chat
}

func (thisUI *UI) setImagesIDs(images []string) {
	thisUI.RWMutex.Lock()
	defer thisUI.RWMutex.Unlock()
	thisUI.imagesIDs = images
}

func (thisUI *UI) getImagesIDs() []string {
	thisUI.RLock()
	defer thisUI.RWMutex.RUnlock()
	return thisUI.imagesIDs
}

func (thisUI *UI) loadChat(chatFrame *cview.Frame, chatView *cview.TextView, chat whatsapp.Chat) func() {
	return func() {

		thisUI.setSelectedContact(chat)
		thisUI.setImagesIDs([]string{})

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

		thisUI.setSelectedContactMessages(messages)

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

func (thisUI *UI) getImageAnsi(msg whatsapp.ImageMessage) string {
	imageScale := 2  //this can be one of 0 - resize (default) or 1 - fill or   2 - fit
	imageDither := 1 //this can be 0 - no dithering (default) or  1 - with blocks or   2 - with chars
	// set image scale factor for ANSIPixel grid
	sfy, sfx := ansimage.BlockSizeY, ansimage.BlockSizeX // 8x4 --> with dithering
	if ansimage.DitheringMode(imageDither) == ansimage.NoDithering {
		sfy, sfx = 2, 1 // 2x1 --> without dithering
	}

	txt := ""
	var contentBytes []byte
	var err error

	if thisUI.ChatView != nil {

		log.Debug("Loading image and converting into an ANSI characters array")

		if client.ImageExists(msg.Info.Id) {
			log.Debugf("Image exists. Loading from disk")
			contentBytes, err = client.ReadImage(msg.Info.Id)
			if err != nil {
				log.Warn("Unable to load image file. Re-downloading..")
			}
			contentBytes, err = msg.Download()
		} else {
			log.Debugf("Image does not exis. Downloading...")
			contentBytes, err = msg.Download()
		}

		if err != nil {
			txt = fmt.Sprintf("Error while downloading image: %s", err)
		} else {
			_, _, width, height := thisUI.ChatView.Box.GetRect()

			err = client.SaveImage(msg, contentBytes)
			if err != nil {
				log.Warnf("Error while saving image: %s", err)
			}
			sm := ansimage.ScaleMode(imageScale)
			dm := ansimage.DitheringMode(1)
			content := bytes.NewReader(contentBytes)
			pix, err := ansimage.NewScaledFromReader(content, height*sfy, width*sfx, color.Black, sm, dm)
			runtime.GOMAXPROCS(runtime.NumCPU())
			pix.SetMaxProcs(runtime.NumCPU())

			if err != nil {
				txt = fmt.Sprintf("Error while rendering the image: %s", err)
			} else {
				log.Debug("Rendering image...")
				txt = pix.RenderExt(false, true)
			}
		}
	} else {
		log.Debug("Chatview not ready. Image won't be displayed")
		txt = "Chatview not ready.. unable to display image"
	}

	return txt

}

func (thisUI *UI) searchMessage(msgID string) interface{} {

	if len(thisUI.getSelectedContactMessages()) == 0 {
		return nil
	}

	for _, msg := range thisUI.getSelectedContactMessages() {
		r := reflect.ValueOf(msg)
		f := reflect.Indirect(r).FieldByName("Info") //Big assumption. Every message, regardles the type, will have a field Info of type MessageInfo
		//TODO add some error checking just in case..
		msgInterface := f.Interface()
		msgInfo := msgInterface.(whatsapp.MessageInfo)
		if msgInfo.Id == msgID {
			return msg
		}
	}

	return nil

}

//Function called when the user presses a key in the chat window
func (thisUI *UI) selectMessage(chatView *cview.TextView) func(tcell.Key) {
	return func(key tcell.Key) {
		log.Info("Call to selectMessage")
		currentSelection := chatView.GetHighlights()

		if key == tcell.KeyEnter {
			if len(currentSelection) > 0 {
				msgID := currentSelection[0]
				chatView.Highlight() //Here I should show the image
				log.Infof("Getting image for message %s", msgID)
				txt := ""

				msg := thisUI.searchMessage(msgID)
				if msg == nil {
					txt = "Message not found"
				} else {
					txt = thisUI.getImageAnsi(msg.(whatsapp.ImageMessage))
				}

				thisUI.MessageModal.SetText(fmt.Sprintf("%s", txt))

				thisUI.Pages.SwitchToPage("modal-message")
			} else {
				chatView.Highlight(thisUI.getImagesIDs()[len(thisUI.getImagesIDs())-1]).ScrollToHighlight() //This happens when I press enter the first time. Scroll to the latest image
			}
		} else if len(currentSelection) > 0 {
			index := 0
			for idx, id := range thisUI.getImagesIDs() {
				if id == currentSelection[0] {
					index = idx
					break
				}
			}

			if key == tcell.KeyTab {
				index = (index - 1 + len(thisUI.getImagesIDs())) % len(thisUI.getImagesIDs())
			} else if key == tcell.KeyBacktab {
				index = (index + 1) % len(thisUI.getImagesIDs())
			} else {
				return
			}
			chatView.Highlight(thisUI.getImagesIDs()[index]).ScrollToHighlight()
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
