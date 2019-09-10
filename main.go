package main

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	cmd "github.com/tech-nico/whatsapp-cli/cmd"
)

func initLogs() {
	log.SetLevel(log.TraceLevel)
	log.SetFormatter(&log.TextFormatter{})
	log.SetReportCaller(true)
}

func main() {
	initLogs()
	cmd.Execute()

	defer fmt.Println("Bye!")

	// t := prompt.Input("> ", completer)
	// //Check prompt.New as in https://github.com/c-bata/kube-prompt/blob/master/main.go#L33
	// fmt.Println("You selected " + t)

}
