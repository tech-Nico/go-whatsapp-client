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
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tech-nico/whatsapp-cli/client"
)

// getCmd represents the get command
var getHistoryCmd = &cobra.Command{
	Use:   "history",
	Short: "Retrieve the chat history",
	Long:  `Retrieve all chat histories or a specific chat history. Defaults to single history. Specify -a to get all histories`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Debug("Call to historyCmd")
		client, err := client.NewClient()
		if err != nil {
			log.Errorf("Error while initiating a new Whatsapp Client: %s", err)
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
	// getHistoryCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
