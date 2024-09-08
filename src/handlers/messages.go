package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"lh-whatsapp/src/database"
	"net/http"
)

// GetMessages Handle /chat endpoint
func GetMessages(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	lhid := vars["lhid"]
	chatId := vars["chatID"]

	messages, err := database.FetchChat(lhid, chatId)
	if err != nil {
		http.Error(w, "Unable to fetch chat", http.StatusInternalServerError)
		return
	}

	// Return the data as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}
