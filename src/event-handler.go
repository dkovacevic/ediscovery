// event_handler.go
package main

import (
	"encoding/json"
	"fmt"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
	"strings"

	"go.mau.fi/whatsmeow/types/events"
)

// Log represents the top-level log structure.
type Log struct {
	Legalhold Kibana
}

// Kibana Legalhold log structure for Kibana
type Kibana struct {
	LHID   string `json:"lhid"`
	ID     string `json:"messageId"`
	ChatID string `json:"chatId"`
	Sent   string `json:"sent"`
	Sender string `json:"sender"`
	Text   string `json:"text"`
	From   string `json:"from"`
	Type   string `json:"type"`

	Group        string `json:"group"`
	Participants string `json:"participants"`
}

func eventHandler(client *whatsmeow.Client, evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		// Initializing the Kibana object
		kibana := Kibana{
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
	}
}

func trace(kibana Kibana) {
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
