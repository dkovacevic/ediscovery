// main.go
package main

import (
	"ediscovery/src/database"
	"ediscovery/src/handlers"
	"ediscovery/src/meow"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
)

func main() {
	// Initialize the database connection
	err := database.NewDB("device.db")
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	go func() {
		_, err := meow.InitWhatsAppClients()
		if err != nil {
			log.Fatalf("Failed to Init WhatsApp clients: %v", err)
		}
	}()

	router := mux.NewRouter()

	// Define the API routes
	router.HandleFunc("/api/signup", handlers.SignUpHandler).Methods("POST")
	router.HandleFunc("/api/login", handlers.LoginHandler).Methods("POST")

	router.Handle("/api/users", handlers.AuthMiddleware(http.HandlerFunc(handlers.GetUsers))).Methods("GET")
	router.Handle("/api/{lhid}/chats", handlers.AuthMiddleware(http.HandlerFunc(handlers.GetChats))).Methods("GET")
	router.Handle("/api/{lhid}/chats/{chatid}/messages", handlers.AuthMiddleware(http.HandlerFunc(handlers.GetMessages))).Methods("GET")
	router.Handle("/api/code", handlers.AuthMiddleware(http.HandlerFunc(handlers.GenerateQRCodeJSON))).Methods("GET")

	// Serve static files from the "./static" directory for the root path "/"
	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./static/"))))

	// Start the server
	fmt.Println("Server started at http://localhost:8080")
	err = http.ListenAndServe(":8080", router)

	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
