// qrcode.go
package main

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/mdp/qrterminal"
	waLog "go.mau.fi/whatsmeow/util/log"
)

func generateQRCode(w http.ResponseWriter, _ *http.Request) {
	deviceStore := container.NewDevice()

	clientLog := waLog.Stdout("Client", "INFO", true)
	client := initializeClient(deviceStore, clientLog)

	qrChan, _ := client.GetQRChannel(context.Background())

	err := client.Connect()
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

				all := strings.ReplaceAll(buf.String(), "\n", "<br>")
				html := fmt.Sprintf(`
					<!DOCTYPE html>
					<html lang="en">
					<head>
						<meta charset="UTF-8">
						<meta name="viewport" content="width=device-width, initial-scale=1.0">
						<title>QR Code</title>
					</head>
					<body>
						<h2>QR Code</h2>
						<pre>%s</pre>
						<h2>Scan this code in your WhatsApp client by openinig Settings, Link Device</h2>
					</body>
					</html>
				`, all)

				w.Header().Set("Content-Type", "text/html")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(html))

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
