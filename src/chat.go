package main

import (
	"encoding/json"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
)

// Handle /chat endpoint
func chatHandler(w http.ResponseWriter, r *http.Request) {
	// Get lhid and chatid from query parameters
	lhid := r.URL.Query().Get("lhid")
	chatid := r.URL.Query().Get("chatid")

	if lhid == "" || chatid == "" {
		http.Error(w, "Missing lhid or chatid", http.StatusBadRequest)
		return
	}

	messages, err := db.fetchChat(lhid, chatid)
	if err != nil {
		http.Error(w, "Unable to fetch chat", http.StatusInternalServerError)
		return
	}

	// Return the data as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}
