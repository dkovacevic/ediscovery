// client.go
package main

import (
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	waLog "go.mau.fi/whatsmeow/util/log"
)

func initializeClient(store *store.Device, logger waLog.Logger) *whatsmeow.Client {
	client := whatsmeow.NewClient(store, logger)
	return client
}
