// client.go
package meow

import (
	"fmt"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
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
		client := NewClient(deviceStore)

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

func NewClient(deviceStore *store.Device) *whatsmeow.Client {
	clientLog := waLog.Stdout("Client", "INFO", true)
	client := whatsmeow.NewClient(deviceStore, clientLog)
	return client
}

func GetUserInfo(client *whatsmeow.Client, jids []types.JID) (map[types.JID]types.UserInfo, error) {
	return client.GetUserInfo(jids)
}

func GetAllDevices() ([]*store.Device, error) {
	return container.GetAllDevices()
}

func NewDevice() *store.Device {
	return container.NewDevice()
}

func GetDevice(jid types.JID) (*store.Device, error) {
	return container.GetDevice(jid)
}
