// event_handler.go
package meow

import (
	"encoding/json"
	"fmt"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
	"lh-whatsapp/src/database"
	"lh-whatsapp/src/models"
	"log"
	"strings"

	"go.mau.fi/whatsmeow/types/events"
)

func EventHandler(client *whatsmeow.Client, evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		// Initializing the Kibana object
		kibana := models.Kibana{
			LHID:   client.Store.ID.User,
			ID:     v.Info.ID,
			ChatID: v.Info.Chat.String(),
			Sent:   v.Info.Timestamp.String(),
			Sender: v.Info.PushName,
			From:   v.Info.Sender.String(),
			Type:   v.Info.Type,
			Text:   v.Message.GetConversation(),
		}

		if v.Info.IsGroup {
			// Fetch group info to get the group title
			groupInfo, err := client.GetGroupInfo(v.Info.Chat)
			if err == nil {
				kibana.Group = groupInfo.Name
				kibana.Participants = extractJIDs(groupInfo.Participants)
			}
		}

		trace(kibana)

		storeDB(kibana)
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

// Custom function to extract JIDs
func extractJIDs(participants []types.GroupParticipant) string {
	if participants == nil {
		return ""
	}

	jids := make([]string, len(participants))
	for i, participant := range participants {
		jids[i] = participant.JID.String()
	}
	return strings.Join(jids, ",")
}
