// event_handler.go
package main

import (
	"encoding/json"
	"fmt"
	"go.mau.fi/whatsmeow"

	"go.mau.fi/whatsmeow/types/events"
)

// Log represents the top-level log structure.
type Log struct {
	Legalhold Kibana
}

// Legalhold log structure for Kibana
type Kibana struct {
	LHID   string `json:"lhid"`
	ID     string `json:"id"`
	Sent   string `json:"sent"`
	Sender string `json:"sender"`
	Text   string `json:"text"`
	From   string `json:"from"`
	To     string `json:"to"`
	Type   string `json:"type"`
}

func eventHandler(client *whatsmeow.Client, evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		// Initializing the Kibana object
		kibana := Kibana{
			LHID:   client.Store.ID.User,
			ID:     v.Info.ID,
			Sent:   v.Info.Timestamp.String(),
			Sender: v.Info.PushName,
			From:   v.Info.Sender.String(),
			Type:   v.Info.Type,
			Text:   v.Message.GetConversation(),
		}

		if v.Info.DeviceSentMeta != nil {
			kibana.To = v.Info.DeviceSentMeta.DestinationJID
		}

		log := Log{
			Legalhold: kibana,
		}

		// Marshal the Kibana object to JSON
		jsonData, err := json.Marshal(log)
		if err != nil {
			fmt.Printf("Error marshaling Kibana object to JSON: %v\n", err)
			return
		}

		// Print the JSON string
		fmt.Println(string(jsonData))
	}
}
