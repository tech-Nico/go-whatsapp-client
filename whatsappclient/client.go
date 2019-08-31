package whatsappclient

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"time"

	qrcodeTerminal "github.com/Baozisoftware/qrcode-terminal-go"
	whatsapp "github.com/Rhymen/go-whatsapp"
	log "github.com/sirupsen/logrus"
)

var loginLogger = log.WithFields(log.Fields{"event": "login", "config_file": getConfigFileName()})

// getConfigFileName Return the full path of the config file based upon the current user's home folder
func getConfigFileName() string {
	home := getHomeFolder()
	return filepath.Join(home, ".go-whatsapp-client/config.conf")
}

func createConfigFileIfNeeded() (*os.File, error) {
	log.Tracef("entered createConfigFile")
	configFileName := getConfigFileName()
	log.Tracef("configFileName: '%s'", configFileName)
	dirStr, _ := path.Split(configFileName)
	log.Tracef("The config folder: %s", dirStr)
	if _, err := os.Stat(configFileName); os.IsNotExist(err) {
		err := os.MkdirAll(dirStr, os.ModePerm)
		if err != nil {
			loginLogger.Errorf("Error while creating folder '%s' : %s", dirStr, err)
		}

		file, err := os.Create(configFileName)

		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Warnf("Error while creating configuration file '%s'", configFileName)
		}

		return file, err
	}

	file, err := os.Open(configFileName)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Warnf("Error while opening the config file '%s'", configFileName)
	}

	return file, err
}

func writeSessionToFile(s whatsapp.Session) error {
	loginLogger.Tracef("Writing session %v to the config file...", s)

	file, err := createConfigFileIfNeeded()
	if err != nil {
		return err
	}

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

// WhatsappClient This is the client object that will allow you to do all necessary actions with your whatsapp account
type WhatsappClient struct {
	Session whatsapp.Session
}

func newLogin(sessionStr string) (whatsapp.Session, error) {
	loginLogger.Tracef("in newLogin with sessionStr '%s'", sessionStr)
	wac, err := whatsapp.NewConn(5 * time.Second)

	if err != nil {
		loginLogger.WithFields(log.Fields{
			"error": err,
		}).Panic("Error while creating a new Whatsapp connection")
	}

	qr := make(chan string)
	if sessionStr == "" {
		log.Debugf("No session passed. Initiate a new login..")
		go func() {
			terminal := qrcodeTerminal.New()
			terminal.Get(<-qr).Print()
		}()
	} else {
		loginLogger.Debugf("Session string.. resuming session.")
		go func() {
			log.Trace("Sending session to channel qr")
			qr <- sessionStr
		}()

	}

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
NewClient Create a new WhatsappClient that will let you do all things with whatsapp.
If a session is stored on disk, use that session otherwise ask to login.
If a session is stored on disk but the session is expired, then ask to login
*/
func NewClient() (WhatsappClient, error) {
	var session whatsapp.Session
	var err error
	configFile := getConfigFileName()
	if FileExists(configFile) {
		loginLogger.Tracef("Config file '%s' exists. Resuming session...", configFile)
		//Try to use the config file as a session
		content, err := ioutil.ReadFile(configFile)

		if err != nil {
			loginLogger.WithFields(log.Fields{
				"error": err,
			})
			loginLogger.Error("Error while trying to open the config file.")
		}
		contentStr := string(content)
		loginLogger.Tracef("Session: %s", contentStr)

		session, err = newLogin(contentStr)
		if err != nil {
			loginLogger.WithFields(log.Fields{
				"error": err,
			}).Error("Error while creating a new Whatsapp session")
		}

	} else {
		loginLogger.Debug("Config file could not be found. Initiating new session...")

		session, err = newLogin("")

		if err != nil {
			loginLogger.WithFields(log.Fields{
				"error": err,
			}).Error("Error while logging in to Whatsapp")
			return WhatsappClient{}, err
		}

		loginLogger.Tracef("Successfully logged in to Whatsapp. Session : %v", session)
		loginLogger.Debug("Storing session to config file")
		err = writeSessionToFile(session)
		if err != nil {
			loginLogger.Warnf("Error while writing config file : %s", err)
		}

	}

	return WhatsappClient{
		Session: session,
	}, nil

}
