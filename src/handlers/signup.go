package handlers

import (
	"ediscovery/src/database"
	"net/http"
)

// SignUpHandler handles the signup form submission
func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse the form data
	username := r.FormValue("username")
	password := r.FormValue("password")

	if password == "" {
		http.Error(w, "Password field is required", http.StatusBadRequest)
		return
	}

	if username != "admin" {
		http.Error(w, "Invalid username", http.StatusBadRequest)
		return
	}

	// Insert the new user into the database
	err := database.InsertUser(username, password)
	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	// Redirect the user to a success page or login page
	http.Redirect(w, r, "/login.html", http.StatusSeeOther)
}
