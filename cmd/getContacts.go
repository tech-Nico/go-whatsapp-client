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

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tech-nico/whatsapp-cli/client"
)

//Whether or not to display all contacts (including anonymous)
var all bool

// getChatsCmd represents the getChats command
var getContactsCmd = &cobra.Command{
	Use:   "contacts",
	Short: "Retrieve the list of chats",
	Long:  `Retrieve the list of chats (1-1 or groups) currently opened`,
	Run:   getContacts,
}

func getContacts(cmd *cobra.Command, args []string) {
	log.Debug("getContacts called")
	wc, err := client.NewClient()
	if err != nil {
		log.Errorf("Error while initializing Whatsapp client: %s", err)
	}
	contacts, err := wc.GetContacts()
	if err != nil {
		log.Fatalf("Error while retriving contacts: %s", err)
	}

	log.Tracef("Contacts: %v", contacts)
	/*	type orderedChat struct {
			Name string
			Time int64
		}
		ordered := make([]orderedChat, 0)
	*/

	if len(contacts) == 0 {
		fmt.Print("No contacts found")
	} else {
		noName := make([]string, 0)
		storedContacts := make([]string, 0)

		for k, v := range contacts {
			if strings.TrimSpace(v.Name) == "" {
				noName = append(noName, strings.TrimSuffix(k, "@s.whatsapp.net"))
			} else {
				storedContacts = append(storedContacts, v.Name+" ("+v.Short+")")
			}

		}

		//Display contacts with no name (just the phone number)
		if all {
			for k := range noName {
				fmt.Printf("%s\n", noName[k])
			}
		}
		//Display stored contacts
		for k := range storedContacts {
			fmt.Printf("%s\n", storedContacts[k])
		}
	}

}

func init() {
	getCmd.AddCommand(getContactsCmd)
	getContactsCmd.Flags().BoolVarP(&all, "all", "a", false, "Display all contacts, including those not stored in your contact list")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getContactsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getContactsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
