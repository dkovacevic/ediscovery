package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"lh-whatsapp/src/database"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func GetChats(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	lhid := vars["lhid"]

	if lhid == "" {
		http.Error(w, "Missing lhid", http.StatusBadRequest)
		return
	}

	chats, err := database.FetchAllChats(lhid)
	if err != nil {
		http.Error(w, "Unable to fetch chats", http.StatusInternalServerError)
		return
	}

	// Return the data as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chats)
}
