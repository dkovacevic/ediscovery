package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mdp/qrterminal"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"net/http"
	"os"
	"strings"
)

var container *sqlstore.Container

func generateQRCode(w http.ResponseWriter, r *http.Request) {
	deviceStore := container.NewDevice()

	clientLog := waLog.Stdout("Client", "INFO", true)
	client := whatsmeow.NewClient(deviceStore, clientLog)

	qrChan, _ := client.GetQRChannel(context.Background())

	err := client.Connect()
	if err != nil {
		panic(err)
	}

	client.AddEventHandler(func(evt interface{}) {
		eventHandler(client, evt)
	})

	var buf bytes.Buffer

	for evt := range qrChan {
		if evt.Event == "code" {
			fmt.Println("QRChannel event: ", evt.Event)

			if buf.Len() == 0 {
				// Create a buffer to capture the QR code output
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, &buf)

				qr := buf.String()

				// Serve HTML with embedded base64 QR code
				all := strings.ReplaceAll(qr, "\n", "<br>")
				html := fmt.Sprintf(`
							<!DOCTYPE html>
							<html lang="en">
							<head>
								<meta charset="UTF-8">
								<meta name="viewport" content="width=device-width, initial-scale=1.0">
								<title>QR Code</title>
							</head>
							<body>
								<h1>QR Code</h1>
								<pre>%s</pre>
								<h2>Scan this code in your WhatsApp client by openinig Settings, Link Device</h2>
								<pre>%s</pre>
							</body>
							</html>
						`, all, all)

				w.Header().Set("Content-Type", "text/html")
				w.WriteHeader(http.StatusOK)
				_, err = w.Write([]byte(html))
			}
		} else {
			fmt.Println("Login event:", evt.Event)
		}
	}
}

func main() {
	dbLog := waLog.Stdout("Database", "INFO", true)

	// Initialize the database connection
	container, _ = sqlstore.New("sqlite3", "file:device.db?_foreign_keys=on", dbLog)

	devices, err := container.GetAllDevices()
	if err != nil {
		panic(err)
	}

	for _, deviceStore := range devices {
		clientLog := waLog.Stdout("Client", "INFO", true)
		client := whatsmeow.NewClient(deviceStore, clientLog)
		client.AddEventHandler(func(evt interface{}) {
			eventHandler(client, evt)
		})

		// Already logged in, just connect
		if client.Store.ID != nil {
			fmt.Println("Connecting WhatsApp client:", client.Store.ID)

			err = client.Connect()
			if err != nil {
				panic(err)
			}
		}
	}

	http.HandleFunc("/link", generateQRCode)

	port := "8080"
	fmt.Println("Server running on port", port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Println("Failed to start server:", err)
		os.Exit(1)
	}

	select {}
}

// Log represents the top-level log structure.
type Log struct {
	Legalhold Kibana
}

// Legalhold log stucture for Kibana
type Kibana struct {
	LHID   string `json:"lhid"`
	ID     string `json:"id"`
	Sent   string `json:"sent"`
	Sender string `json:"sender"`
	Text   string `json:"text"`
	From   string `json:"from"`
	To     string `json:"to"`
	Type   string `json:"type"`
}

func eventHandler(client *whatsmeow.Client, evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		if v.Info.Type == "text" {
			// Initializing the Kibana object
			kibana := Kibana{
				LHID:   client.Store.ID.User,
				ID:     v.Info.ID,
				Sent:   v.Info.Timestamp.String(),
				Sender: v.Info.PushName,
				From:   v.Info.Sender.String(),
				Type:   v.Info.Type,
				Text:   v.Message.GetConversation(),
			}

			if v.Info.DeviceSentMeta != nil {
				kibana.To = v.Info.DeviceSentMeta.DestinationJID
			}

			log := Log{
				Legalhold: kibana,
			}

			// Marshal the Kibana object to JSON
			jsonData, err := json.Marshal(log)
			if err != nil {
				fmt.Printf("Error marshaling Kibana object to JSON: %v\n", err)
				return
			}

			// Print the JSON string
			fmt.Println(string(jsonData))
		}
	}
}
