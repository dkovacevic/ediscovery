package meow

import (
	"context"
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
	var clients []*whatsmeow.Client

	dbLog := waLog.Stdout("Database", "INFO", true)

	container, err = sqlstore.New(context.Background(), "sqlite3", "file:data/device.db?_foreign_keys=on", dbLog)
	if err != nil {
		return nil, err
	}

	devices, err := container.GetAllDevices(context.Background())
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
			EventHandler(deviceStore, evt)
		})

		clients = append(clients, client)
	}
	return clients, nil
}

func NewClient(deviceStore *store.Device) *whatsmeow.Client {
	clientLog := waLog.Stdout("Client", "INFO", true)
	client := whatsmeow.NewClient(deviceStore, clientLog)
	return client
}

func _(client *whatsmeow.Client, jids []types.JID) (map[types.JID]types.UserInfo, error) {
	return client.GetUserInfo(jids)
}

func GetAllDevices() ([]*store.Device, error) {
	return container.GetAllDevices(context.Background())
}

func NewDevice() *store.Device {
	return container.NewDevice()
}

func GetDevice(jid types.JID) (*store.Device, error) {
	return container.GetDevice(context.Background(), jid)
}
