package handlers

import (
	"ediscovery/src/database"
	"ediscovery/src/models"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"go.mau.fi/whatsmeow/types"
	"net/http"
	"strconv"
)

// GetMessages Handle /chat endpoint
func GetMessages(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chatId := vars["chatid"]

	lhJID, err := types.ParseJID(vars["lhid"])

	if err != nil {
		http.Error(w, "Missing lhid", http.StatusBadRequest)
		return
	}

	// Parse pagination parameters
	page, limit := parsePaginationParams(r)

	// Fetch paginated messages from the database
	messages, err := database.FetchPaginatedChat(lhJID.User, chatId, page, limit)
	if err != nil {
		http.Error(w, "Unable to fetch chat", http.StatusInternalServerError)
		fmt.Printf("database.FetchPaginatedChat: %v", err)
		return
	}

	totalMessages, err := database.FetchTotalMessagesCount(lhJID.User, chatId)
	if err != nil {
		http.Error(w, "Unable to fetch total message count", http.StatusInternalServerError)
		fmt.Printf("database.FetchTotalMessagesCount: %v", err)
		return
	}

	response := models.PaginatedResponse{
		Page:       page,
		Limit:      limit,
		Total:      totalMessages,
		TotalPages: (totalMessages + limit - 1) / limit, // Calculate total pages
		Messages:   messages,
	}

	// Return the data as JSON
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}

// Parse pagination parameters from the request, with defaults
func parsePaginationParams(r *http.Request) (int, int) {
	// Default values
	page := 1
	limit := 10

	// Get query parameters
	query := r.URL.Query()

	// Parse 'page'
	if p, err := strconv.Atoi(query.Get("page")); err == nil && p > 0 {
		page = p
	}

	// Parse 'limit'
	if l, err := strconv.Atoi(query.Get("limit")); err == nil && l > 0 {
		limit = l
	}

	return page, limit
}
