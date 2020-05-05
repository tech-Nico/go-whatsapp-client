package client

import (
	"fmt"
	"time"

	qrcodeTerminal "github.com/Baozisoftware/qrcode-terminal-go"
	whatsapp "github.com/Rhymen/go-whatsapp"
	log "github.com/sirupsen/logrus"
)

// WhatsappClient This is the client object that will allow you to do all necessary actions with your whatsapp account
type WhatsappClient struct {
	Session  *whatsapp.Session
	WaC      *whatsapp.Conn
	Chats    []whatsapp.Chat
	Contacts map[string]whatsapp.Contact
}

func newLogin(wac *whatsapp.Conn) error {
	if wac == nil {
		log.Fatal("Whatsapp connection object empty. Please try logging in again.")
	}

	//load saved session
	log.Debug("in newLogin")
	session, err := ReadSession()
	log.WithField("session", session).Trace("session read")
	if err == nil {
		log.Trace("session read successful. Restoring session...")
		//restore session
		session, err = wac.RestoreWithSession(session)
		if err != nil {
			log.WithField("error", err).Warn("error while restoring session. Deleting session file")
			err = deleteSession()
			if err != nil {
				return err
			}
			log.Debug("reattempting login...")
			return newLogin(wac)

		}
		log.WithField("session", session).Trace("session restored")
	} else {
		log.Trace("no saved session -> regular login")
		qr := make(chan string)
		log.Trace("Waiting for qr code to be scanned..")
		go func() {
			terminal := qrcodeTerminal.New()
			terminal.Get(<-qr).Print()
		}()
		log.Trace("Got QR code")
		log.Trace("Attempting login...")
		session, err = wac.Login(qr)

		if err != nil {
			return fmt.Errorf("error during login: %v", err)
		}
	}

	//save session
	log.WithField("session", session).Trace("writing session to file")
	err = WriteSession(session)
	if err != nil {
		return fmt.Errorf("error saving session: %v", err)
	}
	log.Trace("session written to file")
	return nil
}

type IWhatsappHandler interface {
	SetClient(*WhatsappClient)
}

/*
NewClient Create a new WhatsappClient that will let you do all things with whatsapp.
If a session is stored on disk, use that session otherwise ask to login.
If a session is stored on disk but the session is expired, then ask to login
*/
func NewClient(handler interface{}) (WhatsappClient, error) {
	if handler == nil {
		log.Fatalf("Empty handler passed to NewClient")
	}

	if _, ok := handler.(IWhatsappHandler); !ok {
		log.Fatalf("NewClient requires a parameter of type IWhatsappHandler and whatsapp.Handler")
	}

	if _, ok := handler.(whatsapp.Handler); !ok {
		log.Fatalf("NewClient requires a parameter of type IWhatsappHandler and whatsapp.Handler")
	}

	//create new WhatsApp connection
	pwac, err := whatsapp.NewConn(120 * time.Second)

	newClient := WhatsappClient{
		WaC: pwac,
	}

	handler.(IWhatsappHandler).SetClient(&newClient)

	pwac.SetClientName("Command Line Whatsapp Client", "CLI Whatsapp")
	pwac.SetClientVersion(0, 4, 1307)
	if err != nil {
		log.WithField("error", err).Fatal("error creating connection to Whatsapp\n", err)
	}

	pwac.AddHandler(handler.(whatsapp.Handler))

	//login or restore
	if err := newLogin(pwac); err != nil {
		log.WithField("error", err).Fatal("error logging in\n")
	}

	//Disconnect safely
	/*	log.Debug("Shutting down now.")
		session, err := pwac.Disconnect()
		log.Debug("Shut down")
		if err != nil {
			log.WithField("error", err).Fatal("error disconnecting\n")
		}

		log.WithField("session", session).Debug("successfully disconnected from whatsapp")
	*/
	err = newClient.LoadContacts(true)
	if err != nil {
		return WhatsappClient{}, fmt.Errorf("Error initializing WhatsappClient: %s", err)
	}

	//	chats, err := newClient.GetChats()
	err = newClient.GetChats()

	if err != nil {
		return WhatsappClient{}, fmt.Errorf("Error initializing WhatsappClient: %s", err)
	}

	//Wait a couple fo seconds so chats and contacts are downloaded
	<-time.After(time.Second * 2)

	return newClient, nil

}

//RestoreSession Create a new WhatsappClient using the session stored on disk.
//If the session is expired this function won't attempt to login but it will fail and return an error
func RestoreSession() (WhatsappClient, error) {
	//create new WhatsApp connection
	wac, err := whatsapp.NewConn(5 * time.Second)

	if err != nil {
		log.WithField("error", err).Fatal("error creating connection")
	}

	//load saved session
	session, err := ReadSession()
	if err != nil {
		log.WithField("error", err).Fatal("Error while reading the session")
	}

	//restore session
	session, err = wac.RestoreWithSession(session)
	if err != nil {
		log.WithField("error", err).Fatal("Error while restoring the session")
	}

	wc := WhatsappClient{
		Session: &session,
		WaC:     wac,
	}

	return wc, nil

}

// GetContactName returns the name of a contact or the name of a group given its jid number
func (c *WhatsappClient) GetContactName(jid string) string {
	return c.WaC.Store.Contacts[jid].Name
}

//Disconnect terminates whatsapp connection gracefully
func (c *WhatsappClient) Disconnect() error {

	//Disconnect safely
	log.Info("Shutting down now.")
	session, err := c.WaC.Disconnect()
	if err != nil {
		log.WithField("error", err).Fatal("error disconnecting\n")
	}

	log.WithField("session", session).Debug("successfully disconnected from whatsapp")
	err = deleteSession()
	if err != nil {
		return err
	}
	return nil
}
