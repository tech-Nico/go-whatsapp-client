package main

import (
	"fmt"

	go_prompt "github.com/c-bata/go-prompt"
	completer "github.com/c-bata/go-prompt/completer"
	log "github.com/sirupsen/logrus"
	whatsappclient "github.com/tech-nico/go-whatsapp-client/whatsappclient"
)

func initLogs() {
	log.SetLevel(log.TraceLevel)
	log.SetFormatter(&log.TextFormatter{})
	log.SetReportCaller(true)
}

func whatsappCompleter(d go_prompt.Document) []go_prompt.Suggest {
	s := []go_prompt.Suggest{
		{Text: "login", Description: "Login into whatsapp scanning a QR code"},
		{Text: "list chats", Description: "List all the current chats"},
		{Text: "comments", Description: "Store the text commented to articles"},
	}
	return go_prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

var (
	version  string
	revision string
)

func main() {
	comp := whatsappclient.NewCompleter()
	initLogs()
	fmt.Printf("kube-prompt %s (rev-%s)\n", version, revision)
	fmt.Println("Please use `exit` or `Ctrl-D` to exit this program.")
	fmt.Println("Please select table.")
	defer fmt.Println("Bye!")

	p := go_prompt.New(
		whatsappclient.Executor,
		comp.Complete,
		go_prompt.OptionTitle("kube-prompt: interactive kubernetes client"),
		go_prompt.OptionPrefix(">>> "),
		go_prompt.OptionInputTextColor(go_prompt.Yellow),
		go_prompt.OptionCompletionWordSeparator(completer.FilePathCompletionSeparator),
	)
	p.Run()

	// t := prompt.Input("> ", completer)
	// //Check prompt.New as in https://github.com/c-bata/kube-prompt/blob/master/main.go#L33
	// fmt.Println("You selected " + t)

}
