// qrcode.go
package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/mdp/qrterminal"
	waLog "go.mau.fi/whatsmeow/util/log"
	"html/template"
	"io/ioutil"
	"net/http"
)

// QRData Data structure to pass data to the HTML template
type QRData struct {
	QRCode string
}

func generateQRCode(w http.ResponseWriter, _ *http.Request) {
	deviceStore := container.NewDevice()

	clientLog := waLog.Stdout("Client", "INFO", true)
	client := initializeClient(deviceStore, clientLog)

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
		eventHandler(client, evt)
	})

	// Start a goroutine to handle events after sending the QR code
	go func() {
		for evt := range qrChan {
			if evt.Event == "success" {
				// Handle successful connection
				fmt.Println("Login event:", evt.Event)
			}
		}

		// QR was not scanned in time
		if client.Store.ID == nil {
			client.Disconnect()
		}
	}()

	//flusher, _ := w.(http.Flusher)

	// Handle the QR code event in the main thread
	for evt := range qrChan {
		if evt.Event == "code" {
			fmt.Println("QRChannel event: ", evt.Event)

			// Generate and return the QR code immediately
			var buf bytes.Buffer
			qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, &buf)

			// Generate HTML
			renderQR(w, buf)

			//flusher.Flush() // Flush the QR code to the response

			return
		}
	}
}

func renderQR(w http.ResponseWriter, buf bytes.Buffer) {
	data := QRData{
		QRCode: buf.String(),
	}

	// Read the HTML file
	htmlFile, err := ioutil.ReadFile("resources/qr_code.html")
	if err != nil {
		http.Error(w, "Unable to read HTML file", http.StatusInternalServerError)
		return
	}

	// Convert to string
	htmlContent := string(htmlFile)

	// Load the HTML template
	tmpl, err := template.New("qr").Parse(htmlContent)
	if err != nil {
		http.Error(w, "Unable to parse HTML template", http.StatusInternalServerError)
		return
	}

	// Render the HTML with the QR code
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Unable to render template", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
}
