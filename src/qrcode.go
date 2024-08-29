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
	w.Header().Set("Content-Type", "text/html")

	deviceStore := container.NewDevice()

	clientLog := waLog.Stdout("Client", "INFO", true)
	client := initializeClient(deviceStore, clientLog)

	qrChan, err := client.GetQRChannel(context.Background())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = client.Connect()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	client.AddEventHandler(func(evt interface{}) {
		eventHandler(client, evt)
	})

	var buf bytes.Buffer

	for evt := range qrChan {
		fmt.Println("QRChannel event: ", evt.Event)

		if evt.Event == "code" {
			if buf.Len() == 0 {
				// Create a buffer to capture the QR code output
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, &buf)

				if renderQR(w, buf) {
					return
				}

				w.WriteHeader(http.StatusOK)

				flusher, ok := w.(http.Flusher)
				if ok {
					flusher.Flush()
				}
			}
		} else {
			fmt.Println("Login event:", evt.Event)
		}
	}

	if client.Store.ID == nil {
		client.Disconnect()
	}
}

func renderQR(w http.ResponseWriter, buf bytes.Buffer) bool {
	data := QRData{
		QRCode: buf.String(),
	}

	// Read the HTML file
	htmlFile, err := ioutil.ReadFile("resources/qr_code.html")
	if err != nil {
		http.Error(w, "Unable to read HTML file", http.StatusInternalServerError)
		return true
	}

	// Convert to string
	htmlContent := string(htmlFile)

	// Load the HTML template
	tmpl, err := template.New("qr").Parse(htmlContent)
	if err != nil {
		http.Error(w, "Unable to parse HTML template", http.StatusInternalServerError)
		return true
	}

	// Render the HTML with the QR code
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Unable to render template", http.StatusInternalServerError)
	}
	return false
}
