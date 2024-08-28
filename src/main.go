// main.go
package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
	"log"
	"net/http"
	"os"
)

var container *sqlstore.Container
var db *DB

func main() {
	dbLog := waLog.Stdout("Database", "INFO", true)

	// Initialize the database connection
	var err error
	container, err = sqlstore.New("sqlite3", "file:device.db?_foreign_keys=on", dbLog)
	if err != nil {
		panic(err)
	}

	db, err = NewDB("device.db")
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	defer func(Conn *sql.DB) {
		err := Conn.Close()
		if err != nil {
			log.Fatalf("failed to close the database: %v", err)
		}
	}(db.Conn)

	devices, err := container.GetAllDevices()
	if err != nil {
		panic(err)
	}

	for _, deviceStore := range devices {
		clientLog := waLog.Stdout("Client", "INFO", true)
		client := initializeClient(deviceStore, clientLog)

		fmt.Println("Connecting WhatsApp client:", client.Store.ID)

		err = client.Connect()
		if err != nil {
			fmt.Println("Failed to connect the WhatsApp client:", client.Store.ID)
		}

		client.AddEventHandler(func(evt interface{}) {
			eventHandler(client, evt)
		})
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
