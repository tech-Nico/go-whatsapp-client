package whatsappclient

import (
	"encoding/gob"
	"os"
	"path"
	"path/filepath"

	whatsapp "github.com/Rhymen/go-whatsapp"
	log "github.com/sirupsen/logrus"
)

func fileExists(filename string) bool {
	info, err := os.Stat(filename)

	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func getHomeFolder() string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Error while detecting user home folder")
	}

	return home
}

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

func writeSession(session whatsapp.Session) error {
	loginLogger.Tracef("Writing session %v to the config file...", session)
	file, err := createConfigFileIfNeeded()
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := gob.NewEncoder(file)
	err = encoder.Encode(session)
	if err != nil {
		log.Warnf("Error while encoding session: %v", err)
		return err
	}
	return nil
}

func readSession() (whatsapp.Session, error) {
	loginLogger.Debugf("Reading session from file...")
	session := whatsapp.Session{}
	file, err := os.Open(getConfigFileName())
	if err != nil {
		loginLogger.Warnf("Error while opening config file: %v", err)
		return session, err
	}
	defer file.Close()
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&session)
	if err != nil {
		loginLogger.Warnf("Error while decoding session from file: %v", err)
		return session, err
	}
	return session, nil
}
