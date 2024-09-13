package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/mdp/qrterminal"
	"go.mau.fi/whatsmeow"
	waLog "go.mau.fi/whatsmeow/util/log"
	"lh-whatsapp/src/meow"
	"net/http"
)

// QRCodeResponse is the structure for the JSON response
type QRCodeResponse struct {
	QRCode string `json:"qr_code"`
}

// GenerateQRCodeJSON generates the QR code and returns it as JSON
func GenerateQRCodeJSON(w http.ResponseWriter, _ *http.Request) {
	deviceStore := meow.NewDevice()

	clientLog := waLog.Stdout("Client", "INFO", true)
	client := whatsmeow.NewClient(deviceStore, clientLog)

	// Get the QR code event channel, this handles its own timeout internally
	qrChan, err := client.GetQRChannel(context.Background())
	if err != nil {
		http.Error(w, "Failed to get QR channel", http.StatusInternalServerError)
		return
	}

	err = client.Connect()
	if err != nil {
		http.Error(w, "Failed to connect client", http.StatusInternalServerError)
		return
	}

	client.AddEventHandler(func(evt interface{}) {
		meow.EventHandler(client, evt)
	})

	// Start a goroutine to handle events after sending the QR code
	go func() {
		for evt := range qrChan {
			if evt.Event == "success" {
				// Handle successful connection
				fmt.Printf("New LH Device. JID: %v\n", client.Store.ID)
			}
		}

		// QR was not scanned in time
		if client.Store.ID == nil {
			client.Disconnect()
		}
	}()

	// Handle the QR code event in the main thread
	for evt := range qrChan {
		if evt.Event == "code" {
			// Generate the ASCII QR code
			var buf bytes.Buffer
			qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, &buf)

			// Prepare the JSON response
			response := QRCodeResponse{
				QRCode: buf.String(),
			}

			// Send the response as JSON
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
			return
		}
	}
}
