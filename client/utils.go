package client

import (
	"encoding/gob"
	"os"
	"path"
	"path/filepath"
	"time"

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

// getConfigFileName Return the full path of the config file based upon the current user's home folder
func getConfigFileName() string {
	home := getHomeFolder()
	return filepath.Join(home, ".go-whatsapp-client/config.conf")
}

func getChatsFileName() string {
	home := getHomeFolder()
	return filepath.Join(home, ".go-whatsapp-client/chats.bin")
}

func createFileIfNeeded(fileName string) (*os.File, error) {
	var file *os.File
	var err error

	log.Tracef("entered createFileIfNeeded")

	log.Tracef("fileName: '%s'", fileName)
	dirStr, _ := path.Split(fileName)
	log.Tracef("The file folder: %s", dirStr)

	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		err := os.MkdirAll(dirStr, os.ModePerm)
		if err != nil {
			log.Errorf("Error while creating folder '%s' : %s", dirStr, err)
		}

		file, err = os.Create(fileName)

		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Warnf("Error while creating file '%s'", fileName)
		}

	} else {
		if err := os.Remove(fileName); err != nil {
			log.WithField("error", err).Errorf("error while removing file %s", fileName)
		}
		file, err = os.Create(fileName)

		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Warnf("Error while creating file '%s'", fileName)
		}

	}

	return file, err
}

func writeSession(session whatsapp.Session) error {
	log.Tracef("Writing session %v to the config file...", session)
	file, err := createFileIfNeeded(getConfigFileName())
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
	log.Debugf("Reading session from file...")
	session := whatsapp.Session{}
	file, err := os.Open(getConfigFileName())
	if err != nil {
		log.Warnf("Error while opening config file: %v", err)
		return session, err
	}
	defer file.Close()
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&session)
	if err != nil {
		log.Warnf("Error while decoding session from file: %v", err)
		return session, err
	}
	return session, nil
}

func deleteSession() error {
	return os.Remove(getConfigFileName())
}

func storeChatsToFile(chats map[string]Chat) {
	log.Trace("Writing list of chats to file...")
	file, err := createFileIfNeeded(getChatsFileName())
	if err != nil {
		log.Warn("Error while reading chats file: %s. Nothing stored to file", err)
		return
	}
	defer file.Close()
	encoder := gob.NewEncoder(file)
	err = encoder.Encode(chats)
	if err != nil {
		log.Warnf("Error while encoding chats: %v", err)
	}
}

func readChatsFromFile() map[string]Chat {
	log.Debugf("Reading chats from file...")
	chats := make(map[string]Chat)
	//If the file was last updated more than 1 day ago, return nil so new chats will be pulled
	file, err := os.Open(getChatsFileName())
	if err != nil {
		log.Warnf("Error while opening chats file: %v", err)
		return nil
	}
	info, err := os.Stat(getChatsFileName())
	if err != nil {
		log.Warnf("Error while reading chats file stat: %v", err)
		return nil
	}
	fileAge := time.Now().Sub(info.ModTime())
	if fileAge.Hours() > 24 {
		log.Info("File chat exists but older than one day. Pulling new chats..")
		return nil //Return nil so a new list of chats is pulled
	}

	defer file.Close()
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&chats)
	if err != nil {
		log.Warnf("Error while decoding session from file: %v", err)
		return nil
	}
	return chats
}

func formatDate(timestamp uint64) string {
	msgTimestamp := time.Unix(int64(timestamp), 0)
	msgDateYear, msgDateMonth, msgDateDay := msgTimestamp.Date()
	msgDate := time.Date(msgDateYear, msgDateMonth, msgDateDay, 0, 0, 0, 0, time.Local)

	nowYear, nowMonth, nowDay := time.Now().Date()
	nowDate := time.Date(nowYear, nowMonth, nowDay, 0, 0, 0, 0, time.Local)

	dateFmt := "-"
	dateDifference := nowDate.Sub(msgDate).Hours() / 24

	switch dateDifference {
	case 0:
		dateFmt = msgTimestamp.Format("(Today) 3:04:05pm")
	case 1:
		dateFmt = msgTimestamp.Format("(Yesterday) 3:04:05pm")
	case 2, 3, 4:
		dateFmt = msgTimestamp.Format("(" + msgDate.Weekday().String() + ") 3:04:05pm")
	default:
		dateFmt = msgTimestamp.Format("2006-01-02  3:04:05pm")

	}

	return dateFmt
}
