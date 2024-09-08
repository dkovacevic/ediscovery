// client.go
package meow

import (
	"fmt"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
)

var Container *sqlstore.Container

func InitWhatsAppClients() error {
	devices, err := Container.GetAllDevices()
	if err != nil {
		return err
	}

	for _, deviceStore := range devices {
		clientLog := waLog.Stdout("Client", "INFO", true)
		client := whatsmeow.NewClient(deviceStore, clientLog)

		fmt.Println("Connecting WhatsApp client:", client.Store.ID)

		err = client.Connect()
		if err != nil {
			fmt.Println("Failed to connect the WhatsApp client:", client.Store.ID)
			return err
		}

		client.AddEventHandler(func(evt interface{}) {
			EventHandler(client, evt)
		})
	}
	return err
}
