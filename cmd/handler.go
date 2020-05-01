package cmd

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

type CmdHandler struct {
	Incoming chan interface{}
}

func (cmdH *CmdHandler) HandleError(err error) {
	fmt.Println("Error: %s", err)
}

func NewHandler(ch chan interface{}) *CmdHandler {
	if ch == nil {
		log.Fatal("No channel passed to cmd.NewHandler")
	}

	return &CmdHandler{
		Incoming: ch,
	}

}
