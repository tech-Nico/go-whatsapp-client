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
	wc "github.com/tech-nico/whatsapp-cli/client"
)

// logout represents the logout command
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout from whatsapp",
	Long:  `Logout this client from Whatsapp and delete the session stored on disk so that next time you try to run a command you'll be prompted to login`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Call to logoutCmd")
		log.Trace("Attempting to load a session from disk in order to logout that session otherwise not sure what to disconnect from")
		client, err := wc.RestoreSession()

		if err != nil {
			log.WithField("error", err).Fatal("Error while restoring a session in Logout command")
		}
		log.Trace("Session restored from disk")
		client.Disconnect()

		log.Info("Disconnected")

	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
