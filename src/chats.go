package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
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

	chats, err := Db.FetchAllChats(lhid)
	if err != nil {
		http.Error(w, "Unable to fetch chats", http.StatusInternalServerError)
		return
	}

	// Return the data as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chats)
}
