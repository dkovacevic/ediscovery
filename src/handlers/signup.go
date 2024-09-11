package handlers

import (
	"lh-whatsapp/src/database"
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

	if username == "" || password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
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
