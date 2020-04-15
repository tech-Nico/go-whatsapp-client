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

	"github.com/Rhymen/go-whatsapp"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tech-nico/whatsapp-cli/client"
)

var count int

// getCmd represents the get command
var getHistoryCmd = &cobra.Command{
	Use:   "history",
	Short: "Retrieve the chat history",
	Long:  `Retrieve all chat histories or a specific chat history. Defaults to single history. Specify -a to get all histories`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Debug("Call to historyCmd")
		if len(args) == 0 {
			log.Fatal("Please specify a chat name")
		}

		chat := ""
		for k := range args {
			chat = chat + args[k] + " "
		}

		chat = strings.TrimRight(chat, " ")

		log.Infof("Getting history for chat %s", chat)
		foundChats := make(map[string]whatsapp.Chat)

		client, err := client.NewClient()
		if err != nil {
			log.Fatalf("Error while initiating a new Whatsapp Client: %s", err)
		}

		chats, err := client.GetChats()
		if err != nil {
			log.Fatalf("Error while retrieving chats: %s", err)
		}

		//Search all the chats with the given name. Since you might have two chats with the same name, ask which one to use
		//in case there are two or more with the same name
		for k, v := range chats {
			if strings.EqualFold(strings.TrimSpace(v.Name), strings.TrimSpace(chat)) {
				foundChats[k] = v
			}
		}

		selectedJid := ""
		if len(foundChats) == 0 {
			log.Fatalf("There is no such chat '%s'", chat)
		}

		if len(foundChats) == 1 {
			for k, _ := range foundChats {
				selectedJid = k
			}
		}

		if len(foundChats) > 1 {

			fmt.Print("Found the following chats: ")
			for k, v := range foundChats {
				fmt.Printf("[%s] - %s", k, v.Name)
			}
			fmt.Print("Type the Jid to display the chats for")
			fmt.Scan("%s", &selectedJid)

		}

		log.Infof("Get history for jid %s", selectedJid)
		history := client.GetHistory(selectedJid, count)

		for k := range history {
			fmt.Printf("%s", history[k])
		}

		log.Tracef("Logged in to Whatsapp. Session: %v", client.Session)

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
