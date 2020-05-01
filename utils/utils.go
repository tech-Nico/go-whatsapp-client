package utils

import "strings"

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
