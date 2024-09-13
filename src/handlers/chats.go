package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"go.mau.fi/whatsmeow/types"
	"lh-whatsapp/src/database"
	"lh-whatsapp/src/meow"
	"net/http"
)

func GetChats(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	lhJID, err := types.ParseJID(vars["lhid"])

	if err != nil {
		http.Error(w, "Missing lhid", http.StatusBadRequest)
		return
	}

	chats, err := database.FetchAllChats(lhJID.User)
	if err != nil {
		http.Error(w, "Unable to fetch chats", http.StatusInternalServerError)
		return
	}

	device, err := meow.GetDevice(lhJID)

	for i := range chats {
		chatJID, err := types.ParseJID(chats[i].ChatID)
		if err == nil {
			chats[i].PhoneNo = "+" + chatJID.User
			contact, err := device.Contacts.GetContact(chatJID)
			if err == nil {
				if contact.FullName != "" {
					chats[i].Name = contact.FullName
				} else {
					chats[i].Name = contact.PushName
				}
			}
		}
	}

	// Return the data as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chats)
}
