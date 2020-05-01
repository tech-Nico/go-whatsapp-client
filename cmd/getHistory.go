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
	"github.com/Rhymen/go-whatsapp"
	"github.com/Rhymen/go-whatsapp/binary/proto"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var count int

func handleHistoryRawMessage(msg *proto.WebMessageInfo) {
	log.Debug("handleHistoryRawMessage: Handling raw message in getHistory")
	log.Trace(msg)
}

func handleHistoryTextMessage(msg whatsapp.TextMessage) {
	log.Debug("handleHistoryRawMessage: Handling text message in getHistory")
	log.Trace(msg)
}

// getCmd represents the get command
var getHistoryCmd = &cobra.Command{
	Use:   "history",
	Short: "Retrieve the chat history",
	Long:  `Retrieve all chat histories or a specific chat history. Defaults to single history. Specify -a to get all histories`,
	Run: func(cmd *cobra.Command, args []string) {
		/*log.Debug("Call to historyCmd")
		if len(args) == 0 {
			log.Fatal("Please specify a chat name")
		}

		chat := getNameFromArgs(args)

		log.Infof("Getting history for chat %s", chat)

		ch := make(chan interface{})
		wc, err := client.NewClient(ch)

		if err != nil {
			log.Fatalf("Error while initalizing whatsapp client: %s", err)
		}

		chats, err := wc.GetChats()
		if err != nil {
			log.Fatalf("Error while retriving chats: %s", err)
		}

		jid := ""
		for idx := range chats {
			currChat := chats[idx]
			if strings.EqualFold(strings.TrimSpace(currChat.Name), strings.TrimSpace(chat)) {
				jid = currChat.Jid
				break
			}
		}

		if jid == "" {
			log.Fatalf("Unable to find chat %s", chat)
		}

		history, err := wc.GetHistory(jid, count)
		if err != nil {
			log.Fatalf("Error in getting history for chat %s: %s", chat, err)
		}

		for idx := range history {
			fmt.Printf("%s\n", history[idx])
		}
		*/
	},
}

func init() {
	getCmd.AddCommand(getHistoryCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getHistoryCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	getHistoryCmd.Flags().IntVarP(&count, "number", "n", 20, "Number of chat messages to load")
}
