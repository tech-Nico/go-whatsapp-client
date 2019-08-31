package whatsappclient

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func FileExists(filename string) bool {
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
