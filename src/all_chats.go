package main

import (
	"encoding/json"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

// Handle /all-chats endpoint
func allChatsHandler(w http.ResponseWriter, r *http.Request) {
	// Get lhid from query parameters
	lhid := r.URL.Query().Get("lhid")

	if lhid == "" {
		http.Error(w, "Missing lhid", http.StatusBadRequest)
		return
	}

	chats, err := db.fetchAllChats(lhid)
	if err != nil {
		http.Error(w, "Unable to fetch chats", http.StatusInternalServerError)
		return
	}

	// Return the data as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chats)
}
