package main

import (
	log "github.com/sirupsen/logrus"
	wc "github.com/tech-nico/whatsapp-cli/client"
	cmd "github.com/tech-nico/whatsapp-cli/cmd"
)

func initLogs() {
	log.SetFormatter(&log.TextFormatter{
		DisableTimestamp:       false,
		FullTimestamp:          true,
		DisableLevelTruncation: true,
	})
	log.SetReportCaller(false)
}

func main() {
	initLogs()
	c, err := wc.NewClient()
	if err != nil {
		log.Fatalf("Error while creating a new whatsapp client: %s", err)
	}

	log.Debugf("created new whatsapp client! %v", c)
	cmd.Execute()

	// t := prompt.Input("> ", completer)
	// //Check prompt.New as in https://github.com/c-bata/kube-prompt/blob/master/main.go#L33
	// fmt.Println("You selected " + t)

}
