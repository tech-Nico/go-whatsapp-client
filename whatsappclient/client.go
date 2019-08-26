package whatsappclient

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"

	qrcodeTerminal "github.com/Baozisoftware/qrcode-terminal-go"
	whatsapp "github.com/Rhymen/go-whatsapp"
	log "github.com/sirupsen/logrus"
)

const (
	configFile = "~/.go-whatsapp-client/config"
)

var loginLogger = log.WithFields(log.Fields{"event": "login", "config_file": configFile})

// WhatsappClient This is the client object that will allow you to do all necessary actions with your whatsapp account
type WhatsappClient struct {
	Session whatsapp.Session
}

func writeSessionToFile(s whatsapp.Session) error {
	file, err := os.Create(configFile)
	data, _ := json.Marshal(s)

	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.WriteString(file, string(data[:]))
	if err != nil {
		return err
	}

	return file.Sync()
}

func newLogin() (whatsapp.Session, error) {

	wac, err := whatsapp.NewConn(5 * time.Second)

	if err != nil {
		loginLogger.WithFields(log.Fields{
			"error": err,
		}).Panic("Error while creating a new Whatsapp connection")
	}

	qr := make(chan string)
	go func() {
		terminal := qrcodeTerminal.New()
		terminal.Get(<-qr).Print()
	}()

	session, err := wac.Login(qr)
	if err != nil {
		loginLogger.WithFields(log.Fields{
			"error": err,
		}).Error("Error while logging in to Whatsapp")
		return whatsapp.Session{}, err
	}

	loginLogger.Info("Successfully logged in to Whatsapp")
	return session, nil
}

/*
NewClient New create a new WhatsappClient that will allow to do all things with whatsapp.
If a session is stored on disk, use that session otherwise ask to login.
If a session is stored on disk but the session is expired, then ask to login
*/
func NewClient() (WhatsappClient, error) {
	if FileExists(configFile) {
		//Try to use the config file as a session
		content, err := ioutil.ReadFile(configFile)
		if err != nil {
			loginLogger.WithFields(log.Fields{
				"error": err,
			})
			loginLogger.Error("Error while trying to open the config file. Initiating a new session")
		}

		fmt.Printf("File contents: %s", content)
	} else {
		loginLogger.Debug("Config file could not be found. Initiating new session...")
		s, err := newLogin()

		if err != nil {
			loginLogger.WithFields(log.Fields{
				"error": err,
			}).Error("Error while logging in to Whatsapp")
			return WhatsappClient{}, err
		}

		loginLogger.Tracef("Successfully logged in to Whatsapp. Session : %v", s)
		loginLogger.Debug("Storing session to config file")
		err = writeSessionToFile(s)
		if err != nil {
			loginLogger.Errorf("Error while writing config file : %s", err)
		}

	}

	return WhatsappClient{}, nil
}
