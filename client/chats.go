package client

import (
	"sort"
	"strconv"

	whatsapp "github.com/Rhymen/go-whatsapp"
	log "github.com/sirupsen/logrus"
)

func (wc *WhatsappClient) GetChats() error {
	log.Debug("In WhastappClient.GetChats")
	_, err := wc.WaC.Chats()

	if err != nil {
		log.Errorf("Error while retriving chats: %s", err)
		return err
	}

	chats := wc.WaC.Store.Chats
	if err != nil {
		log.Errorf("Error while retrieving chats: %s", err)
		return err
	}

	ordered := make([]whatsapp.Chat, 0)

	log.Debugf("Chats is %v", chats)
	for _, v := range chats {
		ordered = append(ordered, v)
	}

	sort.Slice(ordered, func(i, j int) bool {
		numI, err := strconv.Atoi(ordered[i].LastMessageTime)
		if err != nil {
			log.Errorf("error while converting timestamp %s to number: %s", ordered[i], err)
		}

		numJ, err := strconv.Atoi(ordered[j].LastMessageTime)
		if err != nil {
			log.Errorf("error while converting timestamp %s to number: %s", ordered[j], err)
		}

		return numJ < numI
	})

	wc.Chats = ordered

	return nil
}
