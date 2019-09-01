package whatsappclient

import (
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

func newLogin(existingSession whatsapp.Session) (whatsapp.Session, whatsapp.Conn, error) {
	var session whatsapp.Session

	loginLogger.Tracef("in newLogin with sessionStr '%s'", existingSession)
	wac, err := whatsapp.NewConn(5 * time.Second)
	//Add handler
	handler := &waHandler{wac, make(map[string]struct{})}
	wac.AddHandler(handler)

	if err != nil {
		loginLogger.WithFields(log.Fields{
			"error": err,
		}).Panic("Error while creating a new Whatsapp connection")
	}

	if existingSession.ClientId == "" {
		qr := make(chan string)
		log.Debugf("No session passed. Initiate a new login..")
		go func() {
			terminal := qrcodeTerminal.New()
			terminal.Get(<-qr).Print()
		}()

		session, err = wac.Login(qr)
		if err != nil {
			loginLogger.WithFields(log.Fields{
				"error": err,
			}).Error("Error while logging in to Whatsapp")
			return whatsapp.Session{}, whatsapp.Conn{}, err
		}

		loginLogger.Info("Successfully logged in to Whatsapp")

	} else {
		session, err = wac.RestoreWithSession(existingSession)
		if err != nil {
			loginLogger.WithFields(log.Fields{
				"error": err,
			}).Error("Error while restoring session. Re-login")

			return whatsapp.Session{}, whatsapp.Conn{}, err
		}

	}
	// wait while chat jids are acquired through incoming initial messages
	log.Debug("Waiting for chats info...")
	<-time.After(5 * time.Second)
	log.Debug("Waited for 5 seconds...")
	return session, *wac, nil
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
	var session whatsapp.Session
	var wac whatsapp.Conn
	var err error
	configFile := getConfigFileName()

	if fileExists(configFile) {
		loginLogger.Tracef("Config file '%s' exists. Resuming session...", configFile)
		//Try to use the config file as a session
		session, err = readSession()
		if err != nil {
			loginLogger.WithField("error", err).Error("Error while reading session from config file.")
		}

		session, wac, err = newLogin(session)
		if err != nil {
			loginLogger.WithFields(log.Fields{
				"error": err,
			}).Error("Error while creating a new Whatsapp session")
		}

	} else {
		loginLogger.Debug("Config file could not be found. Initiating new session...")

		session, wac, err = newLogin(whatsapp.Session{})

		if err != nil {
			loginLogger.WithFields(log.Fields{
				"error": err,
			}).Error("Error while logging in to Whatsapp")
			return WhatsappClient{}, err
		}

		loginLogger.Tracef("Successfully logged in to Whatsapp. Session : %v", session)
		loginLogger.Debug("Storing session to config file")
		err = writeSession(session)
		if err != nil {
			loginLogger.Warnf("Error while writing config file : %s", err)
		}

	}

	return WhatsappClient{
		Session: session,
		wac:     wac,
	}, nil

}
