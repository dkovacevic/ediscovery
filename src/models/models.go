package models

type ChatMessage struct {
	SenderName string `json:"sender"`
	Text       string `json:"text"`
	SentDate   string `json:"date"`
}

// Chat represents a single chat overview
type Chat struct {
	ChatID       string `json:"chatId"`
	Name         string `json:"name"`
	GroupName    string `json:"groupName"`
	Participants string `json:"participants"`
	PhoneNo      string `json:"phoneNo"`
}

type ChatListResult struct {
	JID  string `json:"jid"`
	User string `json:"user"`
	Name string `json:"name"`

	Chats []Chat `json:"chats"`
}

type User struct {
	JID    string `json:"jid"`
	User   string `json:"user"`
	Name   string `json:"name"`
	Device string `json:"device"`
}

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
	IsGroup      bool   `json:"is_group"`
}

type PaginatedResponse struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`

	Messages []ChatMessage `json:"messages"`
}
