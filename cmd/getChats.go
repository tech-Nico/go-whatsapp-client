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

	"github.com/Rhymen/go-whatsapp"
	"github.com/spf13/cobra"
	"github.com/tech-nico/whatsapp-cli/utils"
)

var chatSearchStr string

// getChatsCmd represents the getChats command
var getChatsCmd = &cobra.Command{
	Use:   "chats",
	Short: "Retrieve the list of chats",
	Long:  `Retrieve the list of chats (1-1 or groups) currently opened`,
	Run:   getChats,
}

func displayChatSearchResults(searchStr string, chats []whatsapp.Chat) {
	names := make([]string, 100)

	for idx := range chats {
		v := chats[idx]
		name := v.Name
		if v.Name == "" {
			name = v.Jid
		}
		names = append(names, name)
	}
	names = removeDuplicates(names)
	sort.Strings(names)
	names = utils.FilterByContain(names, searchStr)
	fmt.Printf("Chat matching search string '%s':\n%s\n", searchStr, strings.Join(names, "\n"))

}

func displayAllChats(chats []whatsapp.Chat) {

	for k := range chats {
		//time.Unix(ordered[k].Time, 0)
		fmt.Printf("%s\n", chats[k].Name)
	}

}

func getChats(cmd *cobra.Command, args []string) {
	/*	ch := make(chan interface{})
		wc, err := client.NewClient(ch)

		if err != nil {
			log.Fatalf("Error while initializing whatsapp client: %s", err)
		}

		log.Debug("getChats called")

		chats, err := wc.GetChats()
		if err != nil {
			log.Fatalf("Error while getting chats: %s", err)
		}

		if chatSearchStr != "" {
			displayChatSearchResults(chatSearchStr, chats)
		} else {
			displayAllChats(chats)
		}
	*/
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
