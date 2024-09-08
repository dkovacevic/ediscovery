// client.go
package meow

import (
	"fmt"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
)

var container *sqlstore.Container

func InitWhatsAppClients() ([]*whatsmeow.Client, error) {
	var err error
	var clients []*whatsmeow.Client // Slice to store all clients

	dbLog := waLog.Stdout("Database", "INFO", true)

	container, err = sqlstore.New("sqlite3", "file:device.db?_foreign_keys=on", dbLog)

	devices, err := container.GetAllDevices()
	if err != nil {
		return nil, err
	}

	for _, deviceStore := range devices {
		clientLog := waLog.Stdout("Client", "INFO", true)
		client := whatsmeow.NewClient(deviceStore, clientLog)

		fmt.Println("Connecting WhatsApp client:", client.Store.ID)

		err = client.Connect()
		if err != nil {
			fmt.Println("Failed to connect the WhatsApp client:", client.Store.ID)
			return nil, err
		}

		client.AddEventHandler(func(evt interface{}) {
			EventHandler(client, evt)
		})
	}
	return clients, nil
}

func GetAllDevices() ([]*store.Device, error) {
	return container.GetAllDevices()
}

func NewDevice() *store.Device {
	return container.NewDevice()
}
