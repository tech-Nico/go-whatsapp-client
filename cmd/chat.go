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

// chatCmd Open a chat with someone/group
var chatCmd = &cobra.Command{
	Use:   "chats",
	Short: "Retrieve the list of chats",
	Long:  `Retrieve the list of chats (1-1 or groups) currently opened`,
	Run:   chatWith,
}

func chatWith(cmd *cobra.Command, args []string) {
	log.Debug("chatWith called")
	wc, err := client.NewClient()
	if err != nil {
		log.Errorf("Error while initializing Whatsapp client: %s", err)
	}

}

func init() {
	chatCmd.AddCommand(rootCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getChatsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getChatsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}