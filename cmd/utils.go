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

func Filter(vs []string, f func(string) bool) []string {
	vsf := make([]string, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

func FilterByContain(vs []string, searchStr string) []string {
	return Filter(vs, func(a string) bool { return strings.Contains(strings.ToLower(a), strings.ToLower(searchStr)) })
}
