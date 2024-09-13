// event_handler.go
package meow

import (
	"encoding/json"
	"fmt"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/types/events"
	"lh-whatsapp/src/database"
	"lh-whatsapp/src/models"
	"log"
)

func EventHandler(device *store.Device, evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		// Initializing the Kibana object
		kibana := models.Kibana{
			LHID:    device.ID.User,
			ID:      v.Info.ID,
			ChatID:  v.Info.Chat.String(),
			Sent:    v.Info.Timestamp.String(),
			Sender:  v.Info.PushName,
			From:    v.Info.Sender.String(),
			Type:    v.Info.Type,
			Text:    v.Message.GetConversation(),
			IsGroup: v.Info.IsGroup,
		}

		if kibana.Text != "" {
			trace(kibana)

			storeDB(kibana)
		}
	}
}

func storeDB(kibana models.Kibana) {
	err := database.InsertKibana(kibana)
	if err != nil {
		log.Fatalf("failed to insert kibana record: %v", err)
	}
}

func trace(kibana models.Kibana) {
	l := models.Log{
		Legalhold: kibana,
	}

	// Marshal the Kibana object to JSON
	jsonData, err := json.Marshal(l)
	if err != nil {
		fmt.Printf("Error marshaling Kibana object to JSON: %v\n", err)
		return
	}

	// Print the JSON string
	fmt.Println(string(jsonData))
}
