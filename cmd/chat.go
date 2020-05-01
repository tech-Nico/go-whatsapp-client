/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Rhymen/go-whatsapp"
	"github.com/Rhymen/go-whatsapp/binary/proto"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tech-nico/whatsapp-cli/client"
)

// chatCmd Open a chat with someone/group
var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Start a chat",
	Long:  `Start a chat with someone or within an existing group given its name.`,
	Run:   chatWith,
}

func findChat(wc *client.WhatsappClient, chatName string) (whatsapp.Contact, error) {
	allContacts := wc.GetFullContactsDetails(true)

	var foundContact whatsapp.Contact

	for _, v := range allContacts {
		if strings.ToLower(strings.TrimSpace(v.Name)) == strings.ToLower(chatName) {
			foundContact = v
		}
	}

	return foundContact, nil

}

func handleChatRawMessage(msg *proto.WebMessageInfo) {
	log.Debug("chat.handleChatRawMessage: Message from %s", msg.GetKey().GetRemoteJid())
	log.Trace("%v", msg)
}

func handleChatTextMessage(remoteMsgReceived whatsapp.TextMessage, contact whatsapp.Contact) bool {
	log.Debugf("Received message from " + remoteMsgReceived.Info.RemoteJid)
	messageHandled := false
	if remoteMsgReceived.Info.RemoteJid == contact.Jid {
		fmt.Printf("%v\n\t%v\n",
			client.FormatDate(remoteMsgReceived.Info.Timestamp),
			remoteMsgReceived.Text)
		messageHandled = true
	} else {
		log.Debugf("but it's not a message from '%s'... ignoring", contact.Name)
	}

	return messageHandled
}

func sendToChat(ch chan string) {
	for {
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		// convert CRLF to LF
		text = strings.Replace(text, "\n", "", -1)
		ch <- text
	}
}

func chatWith(cmd *cobra.Command, args []string) {

	h := NewHandler(make(chan interface{}))
	wc, err := client.NewClient(h)
	if err != nil {
		log.Errorf("Error while initializing Whatsapp client: %s", err)
	}

	log.Debug("chatWith called")
	if len(args) == 0 {
		log.Fatal("Please specify a contact/group chat")
	}

	chatName := getNameFromArgs(args)
	log.Debugf("Chatting with %s", chatName)
	contact, err := findChat(&wc, chatName)
	if err != nil {
		log.Fatalf("Error while find contact %s", chatName)
	}
	if (whatsapp.Contact{}) == contact {
		log.Fatalf("Contact '%s' could not be found", chatName)
	}
	log.Infof("Contact %s found: %v", chatName, contact)

	//Channel we'll use to send a message to the contact
	sendToChatCh := make(chan string)

	go sendToChat(sendToChatCh)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	fmt.Print("Type ctrl-c to quit.")
	fmt.Print("\n-> ")
ForLoop:
	for {
		select {
		case sendMyMessage := <-sendToChatCh:
			fmt.Println("Just Echoing for now: " + sendMyMessage)
			fmt.Print("\n-> ")
		case remoteMsgReceived := <-h.Incoming:
			switch remoteMsgReceived := remoteMsgReceived.(type) {
			case *proto.WebMessageInfo:
				log.Debugf("chatWith: received raw message: %s", remoteMsgReceived)
			case whatsapp.TextMessage:
				if handleChatTextMessage(remoteMsgReceived, contact) {
					fmt.Print("\n-> ")
				}
			default:
				fmt.Printf("Unknown message type: %v", remoteMsgReceived)
			}
		case interrupted := <-c:
			fmt.Println()
			fmt.Println(interrupted)
			//Disconnect safe
			/*			fmt.Println("Shutting down now.")
						session, err := wc.WaC.Disconnect()
						if err != nil {
							log.Fatalf("error disconnecting: %v\n", err)
						}
						if err := client.WriteSession(session); err != nil {
							log.Fatalf("error saving session: %v", err)
						}*/
			break ForLoop
		}
	}

}

func init() {
	rootCmd.AddCommand(chatCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getChatsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getChatsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
