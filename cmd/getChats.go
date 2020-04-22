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
	"fmt"
	"strings"

	"sort"
	"strconv"

	"github.com/Rhymen/go-whatsapp"
	"github.com/Rhymen/go-whatsapp/binary/proto"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tech-nico/whatsapp-cli/client"
)

var chatSearchStr string

// getChatsCmd represents the getChats command
var getChatsCmd = &cobra.Command{
	Use:   "chats",
	Short: "Retrieve the list of chats",
	Long:  `Retrieve the list of chats (1-1 or groups) currently opened`,
	Run:   getChats,
}

func handleRawMessage(msg *proto.WebMessageInfo) {
	log.Debug("getChats.handleRawMessage invoked. Doing nothing")
	log.Trace(msg)
}

func handleTextMessage(msg whatsapp.TextMessage) {
	log.Debug("getChats.handleTextMessage invoked. Doing nothing")
	log.Trace(msg)
}

func displayChatSearchResults(searchStr string, chats map[string]whatsapp.Chat) {
	names := make([]string, 100)

	for _, v := range chats {
		name := v.Name
		if v.Name == "" {
			name = v.Jid
		}
		names = append(names, name)
	}
	names = removeDuplicates(names)
	sort.Strings(names)
	names = FilterByContain(names, searchStr)
	fmt.Printf("Chat matching search string '%s':\n%s\n", searchStr, strings.Join(names, "\n"))

}

func displayAllChats(chats map[string]whatsapp.Chat) {
	type orderedChat struct {
		Name string
		Time int64
	}

	ordered := make([]orderedChat, 0)

	log.Debugf("Chats is %v", chats)
	for _, v := range chats {

		time, err := strconv.ParseInt(v.LastMessageTime, 10, 64)
		if err != nil {
			log.Warnf("Error while converting a timestamp to integer: %s. Skipping this chat...", err)
			continue
		}
		name := v.Name
		if name == "" {
			name = v.Jid
		}

		newChat := orderedChat{
			Name: name,
			Time: time,
		}

		ordered = append(ordered, newChat)
	}

	sort.Slice(ordered, func(i, j int) bool {
		return ordered[i].Time < ordered[j].Time
	})

	for k := range ordered {
		//time.Unix(ordered[k].Time, 0)
		fmt.Printf("%s\n", ordered[k].Name)
	}

}

func goGetChats(ch chan interface{}) {
	wc, err := client.NewClient(ch)
	if err != nil {
		log.Errorf("Error while initializing Whatsapp client: %s", err)
	}
	chats, err := wc.GetChats()
	if err != nil {
		log.Fatalf("Error while retrieving chats: %s", err)
	}

	ch <- chats
}

func doGetChats() map[string]whatsapp.Chat {
	ch := make(chan interface{})
	chats := map[string]whatsapp.Chat{}
	go goGetChats(ch)
ForLoop:
	for {
		select {
		case msg := <-ch:
			switch msg := msg.(type) {
			case *proto.WebMessageInfo:
				handleRawMessage(msg)
			case whatsapp.TextMessage:
				handleTextMessage(msg)
			case map[string]whatsapp.Chat:
				chats = msg
				break ForLoop
			default:
				log.Warn("\nUnknown message type: %T", msg)
			}

		}
	}

	return chats
}

func getChats(cmd *cobra.Command, args []string) {
	log.Debug("getChats called")

	chats := doGetChats()

	if chatSearchStr != "" {
		displayChatSearchResults(chatSearchStr, chats)
	} else {
		displayAllChats(chats)
	}

}

func init() {
	getCmd.AddCommand(getChatsCmd)
	getChatsCmd.Flags().StringVarP(&chatSearchStr, "search", "s", "", "Search for a chat resembling the given string. Ex: -s twl would return cartwheel")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getChatsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getChatsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
