package client

import (
	"fmt"
	"time"

	qrcodeTerminal "github.com/Baozisoftware/qrcode-terminal-go"
	whatsapp "github.com/Rhymen/go-whatsapp"
	"github.com/Rhymen/go-whatsapp/binary/proto"
	log "github.com/sirupsen/logrus"
)

// WhatsappClient This is the client object that will allow you to do all necessary actions with your whatsapp account
type WhatsappClient struct {
	Session  whatsapp.Session
	wac      whatsapp.Conn
	chats    map[string]Chat
	contacts map[string]whatsapp.Contact
}

type Chat struct {
	Name string
	Jid  string
}

type waHandler struct {
	c     *whatsapp.Conn
	chats map[string]*proto.WebMessageInfo
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
			log.Infof("Got new terminal: %v", terminal)
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
		h.chats[*message.Key.RemoteJid] = message
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
	pwac, err := whatsapp.NewConn(5 * time.Second)

	pwac.SetClientName("Command Line Whatsapp Client", "CLI Whatsapp")
	pwac.SetClientVersion(0, 4, 1307)
	if err != nil {
		log.WithField("error", err).Fatal("error creating connection to Whatsapp\n", err)
	}

	//Add handler
	handler := &waHandler{pwac, make(map[string]*proto.WebMessageInfo)}
	pwac.AddHandler(handler)

	//login or restore
	if err := newLogin(pwac); err != nil {
		log.WithField("error", err).Fatal("error logging in\n")
	}

	var newClient = WhatsappClient{
		wac: *pwac,
	}

	chats := readChatsFromFile()

	if chats == nil {
		//wait while chat jids are acquired through incoming initial messages
		log.Info("Waiting for chats info...")
		<-time.After(5 * time.Second)

		chats := make(map[string]Chat)
		for chatJid := range handler.chats {
			contactName := newClient.GetContactName(chatJid)
			//webMessageInfo := handler.chats[chatJid]
			log.Tracef("Chat: %v: %v", chatJid, contactName)
			newChat := Chat{
				//Name: webMessageInfo.GetMessage().GetContactMessage().GetDisplayName(),
				Name: contactName,
				Jid:  chatJid,
			}
			chats[chatJid] = newChat
		}
		//store chats so we don't have to query all the times
		storeChatsToFile(chats)
	}

	newClient.chats = chats
	for jit, chat := range chats {
		log.Debugf("Jid: %s, Name: %s", jit, chat)
	}

	//Disconnect safely
	/*	log.Debug("Shutting down now.")
		session, err := pwac.Disconnect()
		if err != nil {
			log.WithField("error", err).Fatal("error disconnecting\n")
		}

		log.WithField("session", session).Debug("successfully disconnected from whatsapp")
	*/
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
	session, err := readSession()
	if err != nil {
		log.WithField("error", err).Fatal("Error while reading the session")
	}

	//restore session
	session, err = wac.RestoreWithSession(session)
	if err != nil {
		log.WithField("error", err).Fatal("Error while restoring the session")
	}

	wc := WhatsappClient{
		Session: session,
		wac:     *wac,
	}

	return wc, nil

}

// GetContactName returns the name of a contact or the name of a group given its jid number
func (c *WhatsappClient) GetContactName(jid string) string {
	return c.wac.Store.Contacts[jid].Name
}

func (c *WhatsappClient) GetChats() (map[string]whatsapp.Chat, error) {
	log.Debug("In WhastappClient.GetChats")
	_, err := c.wac.Chats()

	if err != nil {
		log.Errorf("Error while retriving chats: %s", err)
		return nil, err
	}
	return c.wac.Store.Chats, nil
}

//GetContacts return the list of contacts known by Whatsapp
func (c *WhatsappClient) GetContacts() (map[string]whatsapp.Contact, error) {
	log.Debug("In WhatsappClient.GetContacts")
	_, err := c.wac.Contacts()
	if err != nil {
		log.Errorf("Error while retriving contacts: %s", err)
		return nil, err
	}

	return c.wac.Store.Contacts, nil
}

//Disconnect terminates whatsapp connection gracefully
func (c *WhatsappClient) Disconnect() error {

	//Disconnect safely
	log.Info("Shutting down now.")
	session, err := c.wac.Disconnect()
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
