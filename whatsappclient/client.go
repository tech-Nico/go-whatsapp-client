package whatsappclient

import (
	"fmt"
	"time"

	qrcodeTerminal "github.com/Baozisoftware/qrcode-terminal-go"
	whatsapp "github.com/Rhymen/go-whatsapp"
	"github.com/Rhymen/go-whatsapp/binary/proto"
	log "github.com/sirupsen/logrus"
)

var loginLogger = log.WithFields(log.Fields{"event": "login", "config_file": getConfigFileName()})

// getConfigFileName Return the full path of the config file based upon the current user's home folder

// WhatsappClient This is the client object that will allow you to do all necessary actions with your whatsapp account
type WhatsappClient struct {
	Session whatsapp.Session
	wac     whatsapp.Conn
}

type waHandler struct {
	c     *whatsapp.Conn
	chats map[string]struct{}
}

func newLogin(wac *whatsapp.Conn) error {
	//load saved session
	log.Debug("in newLogin")
	session, err := readSession()
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
		session, err = wac.Login(qr)

		if err != nil {
			return fmt.Errorf("error during login: %v", err)
		}
	}

	//save session
	log.WithField("session", session).Trace("writing session to file")
	err = writeSession(session)
	if err != nil {
		return fmt.Errorf("error saving session: %v", err)
	}
	log.Trace("session written to file")
	return nil
}

func (h *waHandler) ShouldCallSynchronously() bool {
	return true
}

func (h *waHandler) HandleRawMessage(message *proto.WebMessageInfo) {
	// gather chats jid info from initial messages
	if message != nil && message.Key.RemoteJid != nil {
		h.chats[*message.Key.RemoteJid] = struct{}{}
	}
}

//HandleError needs to be implemented to be a valid WhatsApp handler
func (h *waHandler) HandleError(err error) {

	if e, ok := err.(*whatsapp.ErrConnectionFailed); ok {
		log.Printf("Connection failed, underlying error: %v", e.Err)
		log.Println("Waiting 30sec...")
		<-time.After(30 * time.Second)
		log.Println("Reconnecting...")
		err := h.c.Restore()
		if err != nil {
			log.Fatalf("Restore failed: %v", err)
		}
	} else {
		log.Printf("error occoured: %v\n", err)
	}
}

/*
NewClient Create a new WhatsappClient that will let you do all things with whatsapp.
If a session is stored on disk, use that session otherwise ask to login.
If a session is stored on disk but the session is expired, then ask to login
*/
func NewClient() (WhatsappClient, error) {
	//create new WhatsApp connection
	wac, err := whatsapp.NewConn(5 * time.Second)
	if err != nil {
		log.WithField("error", err).Fatal("error creating connection to Whatsapp\n", err)
	}

	//Add handler
	handler := &waHandler{wac, make(map[string]struct{})}
	wac.AddHandler(handler)

	//login or restore
	if err := newLogin(wac); err != nil {
		log.WithField("error", err).Fatal("error logging in\n")
	}

	// wait while chat jids are acquired through incoming initial messages
	fmt.Println("Waiting for chats info...")
	<-time.After(5 * time.Second)

	for chat := range handler.chats {
		log.Debugf("Chat: %v: %v", chat, handler.chats[chat])
	}
	//Disconnect safely
	log.Info("Shutting down now.")
	session, err := wac.Disconnect()
	if err != nil {
		log.WithField("error", err).Fatal("error disconnecting\n")
	}

	log.WithField("session", session).Debug("successfully disconnected from whatsapp")

	return WhatsappClient{
		wac: *wac,
	}, nil

}

func (c *WhatsappClient) Disconnect() error {

	// get history synchronously
	//	GetAnyHistory(wac, handler.chats)
	// fmt.Println("Done. Press Ctrl+C for exit.")

	// c := make(chan os.Signal, 1)
	// signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	// <-c

	//Disconnect safely
	log.Info("Shutting down now.")
	session, err := c.wac.Disconnect()
	if err != nil {
		log.WithField("error", err).Fatal("error disconnecting\n")
	}

	log.WithField("session", session).Debug("successfully disconnected from whatsapp")

	return nil

}
