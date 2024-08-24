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
	"time"
)

var container *sqlstore.Container

func generateQRCode(w http.ResponseWriter, r *http.Request) {
	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		panic(err)
	}

	clientLog := waLog.Stdout("Client", "INFO", true)
	client := whatsmeow.NewClient(deviceStore, clientLog)

	if client.Store.ID != nil {
		fmt.Println("DB is not empty. ID: ", client.Store.ID)
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("The Device is already linked"))
		return
	}

	qrChan, _ := client.GetQRChannel(context.Background())

	err = client.Connect()
	if err != nil {
		panic(err)
	}

	client.AddEventHandler(eventHandler)

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
				if err != nil {
					http.Error(w, "Failed to write response", http.StatusInternalServerError)
					return
				}
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

	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		panic(err)
	}

	clientLog := waLog.Stdout("Client", "INFO", true)
	client := whatsmeow.NewClient(deviceStore, clientLog)
	client.AddEventHandler(eventHandler)

	// Already logged in, just connect
	if client.Store.ID != nil {
		fmt.Println("Connecting WhatsApp client:", client.Store.ID)

		err = client.Connect()
		if err != nil {
			panic(err)
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
	ID     string `json:"id"`
	Sent   int64  `json:"sent"`
	Sender string `json:"sender"`
	Text   string `json:"text"`
	From   string `json:"from"`
}

func eventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		// Initializing the Kibana object
		kibana := Kibana{
			ID:     v.Info.ID,
			Sent:   time.Now().Unix(),
			Sender: v.Info.PushName,
			Text:   v.Message.GetConversation(),
			From:   v.Info.Sender.String(),
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
