// main.go
package main

import (
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"lh-whatsapp/src/database"
	"lh-whatsapp/src/handlers"
	"lh-whatsapp/src/meow"
	"log"
	"net/http"
)

func main() {
	// Initialize the database connection
	err := database.NewDB("device.db")
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	_, err = meow.InitWhatsAppClients()
	if err != nil {
		log.Fatalf("Failed to Init WhatsApp clients: %v", err)
	}

	router := mux.NewRouter()

	// Define the API routes
	router.HandleFunc("/api/users", handlers.GetUsers).Methods("GET")
	router.HandleFunc("/api/{lhid}/chats", handlers.GetChats).Methods("GET")
	router.HandleFunc("/api/{lhid}/chats/{chatID}/messages", handlers.GetMessages).Methods("GET")

	// Handle the "/link" route separately
	router.HandleFunc("/link", handlers.GenerateQRCode).Methods("GET")

	// Serve static files from the "./static" directory for the root path "/"
	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./static/"))))

	// Start the server
	fmt.Println("Server started at http://localhost:8080")
	err = http.ListenAndServe(":8080", router)

	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
