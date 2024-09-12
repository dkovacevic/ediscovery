package handlers

import (
	"github.com/golang-jwt/jwt/v5"
	"lh-whatsapp/src/database"
	"net/http"
	"net/url"
	"time"
)

// Claims structure for JWT
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

var jwtKey = []byte("my_secret_key")

const cookieName = "auth5"

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	// Validate input
	if username == "" || password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	// Authenticate the user
	authenticated, err := database.AuthenticateUser(username, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if authenticated {
		// Generate JWT token
		expirationTime := time.Now().Add(15 * time.Minute)
		claims := &Claims{
			Username: username,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expirationTime),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			http.Error(w, "Error creating token", http.StatusInternalServerError)
			return
		}

		// Set JWT token in a cookie
		http.SetCookie(w, &http.Cookie{
			Name:     cookieName,
			Value:    tokenString,
			Expires:  expirationTime,
			HttpOnly: true,
		})

		// Get the redirect URL from the query parameters
		redirectURL := r.FormValue("redirectUrl")
		if redirectURL == "" {
			redirectURL = "/index.html" // Default redirect if none is provided
		}

		// On success, redirect to the original destination or default page
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
	} else {
		http.Error(w, "Invalid login credentials", http.StatusUnauthorized)
	}
}

// AuthMiddleware validates JWT in cookies and redirects to login if token is missing or invalid
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		referrer := r.Referer()
		if referrer == "" {
			referrer = "/" // Default page if no referrer is available
		}

		// Check if the token exists in the cookie
		cookie, err := r.Cookie(cookieName)
		if err != nil {
			if err == http.ErrNoCookie {
				http.Redirect(w, r, "/login.html?redirect="+url.QueryEscape(referrer), http.StatusSeeOther)
				return
			}
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		tokenStr := cookie.Value
		claims := &Claims{}

		// Parse the JWT token and validate
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				http.Redirect(w, r, "/login.html?redirect="+url.QueryEscape(referrer), http.StatusSeeOther)
				return
			}
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		if !token.Valid {
			http.Redirect(w, r, "/login.html?redirect="+url.QueryEscape(referrer), http.StatusSeeOther)
			return
		}

		// Proceed with the request if authenticated
		next.ServeHTTP(w, r)
	})
}
