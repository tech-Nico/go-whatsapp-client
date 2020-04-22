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
	"sort"
	"strings"

	"github.com/Rhymen/go-whatsapp"
	"github.com/Rhymen/go-whatsapp/binary/proto"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tech-nico/whatsapp-cli/client"
)

//Whether or not to display all contacts (including anonymous)
var all bool
var searchStr string

// getChatsCmd represents the getChats command
var getContactsCmd = &cobra.Command{
	Use:   "contacts",
	Short: "Retrieve the list of contacts",
	Long: `Retrieve the list of contacts known by whatsapp. Display the contact ID (number@s.whatsapp.net) if the contact name is empty.
	Use -s if you only know some part of the name and you want to search for a contact containing that part`,
	Run: getContacts,
}

func handleContactsRawMessage(msg *proto.WebMessageInfo) {
	log.Debug("handleContactsRawMessage: Handling raw message in getContacts. Doing nothing...")
	log.Trace(msg)
}

func handleContactsTextMessage(msg whatsapp.TextMessage) {
	log.Debug("handleContactsTextMessage: Handling text message in getContacts. Doing nothing...")
	log.Trace(msg)
}

func goContacts(ch chan interface{}) {

	wc, err := client.NewClient(ch)
	if err != nil {
		log.Errorf("Error while initializing Whatsapp client: %s", err)
	}
	contacts, err := wc.GetContacts()
	if err != nil {
		log.Fatalf("Error while retriving contacts: %s", err)
	}

	log.Tracef("Contacts: %v", contacts)

	ch <- contacts
}

func contactsToStringSlice(contacts map[string]whatsapp.Contact) []string {
	returnContacts := []string{}
	if len(contacts) == 0 {
		log.Info("No contacts found")
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

		if searchStr != "" {
			returnContacts = append(noName, storedContacts...)
			returnContacts = FilterByContain(returnContacts, searchStr)
		} else {
			//order the storedContacts alphabetically
			sort.Strings(storedContacts)

			returnContacts = storedContacts

			//Display contacts with no name (just the phone number)
			if all {
				returnContacts = append(noName, storedContacts...)
			}
		}
	}

	return returnContacts
}

func doGetContacts() map[string]whatsapp.Contact {
	ch := make(chan interface{})
	contacts := make(map[string]whatsapp.Contact, 0)
	go goContacts(ch)
ForLoop:
	for {
		select {
		case msg := <-ch:
			switch msg := msg.(type) {
			case *proto.WebMessageInfo:
				handleContactsRawMessage(msg)
			case whatsapp.TextMessage:
				handleContactsTextMessage(msg)
			case map[string]whatsapp.Contact:
				contacts = msg
				break ForLoop
			default:
				fmt.Printf("Unknown message type %T: %v", msg, msg)
			}
		}
	}

	return contacts
}

func getContacts(cmd *cobra.Command, args []string) {
	log.Debug("getContacts called")
	contacts := contactsToStringSlice(doGetContacts())
	if searchStr != "" {
		if len(contacts) != 0 {
			fmt.Printf("\nMatches found for '%s':\n\n%s\n", searchStr, strings.Join(contacts, "\n"))
		} else {
			fmt.Printf("No contacts containing '%s' was found", searchStr)
		}
	} else {
		//Display stored contacts
		for k := range contacts {
			fmt.Printf("%s\n", contacts[k])
		}

	}

}

func init() {
	getCmd.AddCommand(getContactsCmd)
	getContactsCmd.Flags().BoolVarP(&all, "all", "a", false, "Display all contacts, including those not stored in your contact list")
	getContactsCmd.Flags().StringVarP(&searchStr, "search", "s", "", "Search for a contact resembling the given string. Ex: -s twl would return cartwheel")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getContactsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getContactsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
