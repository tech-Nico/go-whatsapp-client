package cmd

import "strings"

func getNameFromArgs(args []string) string {
	chat := ""
	for k := range args {
		chat = chat + args[k] + " "
	}

	chat = strings.TrimRight(chat, " ")
	return chat
}

func removeDuplicates(names []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range names {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
