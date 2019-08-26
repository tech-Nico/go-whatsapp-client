package whatsapp-client

import (
	"github.com/tech-nico/go-whatsapp-client"
	"github.com/Baozisoftware/qrcode-terminal-go"
	"fmt"
	"os"
	"time"

)

const (
	configFile "~/.go-whatsapp-client/config"
)

type WhatsappClient struct {
	session whatsapp.Session
}

func newLogin(){

	wac, err := whatsapp.NewConn(5 * time.Second)

	if err != nil {
		panic(err)
	}

	qr := make(chan string)
	go func() {
		terminal := qrcodeTerminal.New()
		terminal.Get(<-qr).Print()
	}()

	session, err := wac.Login(qr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error during login: %v\n", err)
	}

	fmt.Printf("login successful, session %v\n", session)

}

/* 
	New create a new WhatsappClient that will allow to do all things with whatsapp.
	If a session is stored on disk, use that session otherwise ask to login.
	If a session is stored on disk but the session is expired, then ask to login
*/
func (client *WhatsappClient) New() error {
	if FileExists(configFile) {
		//Try to use the config file as a session
		
	} else {
	}
}
