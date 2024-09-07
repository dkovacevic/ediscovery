package main

import (
	_ "github.com/mattn/go-sqlite3" // SQLite driver, use the correct driver for your DB
	"html/template"
	"log"
	"net/http"
)

// Handler to generate HTML
func chatHandler(w http.ResponseWriter, r *http.Request) {
	// Extract lhid and chatid from query params
	lhid := r.URL.Query().Get("lhid")
	chatid := r.URL.Query().Get("chatid")

	if lhid == "" || chatid == "" {
		http.Error(w, "Missing lhid or chatid", http.StatusBadRequest)
		return
	}

	// Fetch chat messages
	messages, err := db.fetchChat(lhid, chatid)
	if err != nil {
		http.Error(w, "Unable to fetch chat", http.StatusInternalServerError)
		log.Println("Error fetching chat:", err)
		return
	}

	// Parse and render the HTML template
	tmpl := template.Must(template.ParseFiles("resources/chat.html"))

	err = tmpl.Execute(w, struct {
		ChatID   string
		Messages []ChatMessage
	}{
		ChatID:   chatid,
		Messages: messages,
	})

	if err != nil {
		http.Error(w, "Unable to render template", http.StatusInternalServerError)
		log.Println("Error rendering template:", err)
	}
}
