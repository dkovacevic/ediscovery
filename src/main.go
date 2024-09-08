// main.go
package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
	"log"
	"net/http"
	"os"
)

var Container *sqlstore.Container

func main() {
	dbLog := waLog.Stdout("Database", "INFO", true)

	// Initialize the database connection
	var err error
	Container, err = sqlstore.New("sqlite3", "file:device.db?_foreign_keys=on", dbLog)
	if err != nil {
		panic(err)
	}

	Db, err = NewDB("device.db")
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	defer func(Conn *sql.DB) {
		err := Conn.Close()
		if err != nil {
			log.Fatalf("failed to close the database: %v", err)
		}
	}(Db.Conn)

	devices, err := Container.GetAllDevices()
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

	router := mux.NewRouter()

	// Define the API routes
	router.HandleFunc("/api/users", GetUsers).Methods("GET")
	router.HandleFunc("/api/{lhid}/chats", GetChats).Methods("GET")
	router.HandleFunc("/api/{lhid}/chats/{chatID}/messages", GetMessages).Methods("GET")

	// Handle the "/link" route separately
	router.HandleFunc("/link", generateQRCode).Methods("GET")

	// Serve static files from the "./static" directory for the root path "/"
	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./static/"))))

	// Start the server
	fmt.Println("Server started at http://localhost:8080")
	err = http.ListenAndServe(":8080", router)

	if err != nil {
		fmt.Println("Failed to start server:", err)
		os.Exit(1)
	}
}
