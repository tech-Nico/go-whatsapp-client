package client

import (
	"sort"
	"strings"

	"github.com/Rhymen/go-whatsapp"
	log "github.com/sirupsen/logrus"
	"github.com/tech-nico/whatsapp-cli/utils"
)

func (wc *WhatsappClient) LoadContacts(force bool) error {

	if len(wc.Contacts) > 0 && !force {
		log.Debug("Contacts already loaded and force was false")
		return nil
	}

	log.Debug("In WhatsappClient.GetContacts")

	_, err := wc.WaC.Contacts()
	if err != nil {
		log.Errorf("Error getting whatsapp contacts: %s", err)
		return err
	}

	log.Tracef("Returning contacts %v", wc.WaC.Store.Contacts)

	wc.Contacts = wc.WaC.Store.Contacts

	return nil
}

func (wc *WhatsappClient) GetFullContactsDetails(all bool) map[string]whatsapp.Contact {
	return wc.Contacts
}

func (wc *WhatsappClient) GetFilteredContactsNames(filter string, all bool) ([]string, error) {

	returnContacts, err := wc.GetContactsNames(all)
	if err != nil {
		return []string{}, nil
	}

	if len(returnContacts) == 0 {
		log.Info("No contacts found")
	} else {
		if filter != "" {
			returnContacts = utils.FilterByContain(returnContacts, filter)
		}
	}

	return returnContacts, nil
}

func (wc *WhatsappClient) GetContactsNames(all bool) ([]string, error) {
	err := wc.LoadContacts(false)
	if err != nil {
		return []string{}, err
	}

	returnContacts := []string{}
	if len(wc.Contacts) == 0 {
		log.Info("No contacts found")
	} else {
		noName := make([]string, 0)
		storedContacts := make([]string, 0)

		for k, v := range wc.Contacts {
			if strings.TrimSpace(v.Name) == "" {
				noName = append(noName, strings.TrimSuffix(k, "@s.whatsapp.net"))
			} else {
				storedContacts = append(storedContacts, v.Name+" ("+v.Short+")")
			}

		}

		//order the storedContacts alphabetically
		sort.Strings(storedContacts)
		returnContacts = storedContacts

		//Display contacts with no name (just the phone number)
		if all {
			returnContacts = append(noName, storedContacts...)
		}
	}

	return returnContacts, nil
}
