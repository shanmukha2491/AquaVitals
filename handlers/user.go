package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/shanmukha2491/AquaVitals/config"
	auth "github.com/shanmukha2491/AquaVitals/middlewares"
	"github.com/shanmukha2491/AquaVitals/model"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Error Reading Data", http.StatusPartialContent)
		return
	}

	user, err = config.FindOne(user.Email, user.UserName)
	if err != nil {
		http.Error(w, "User Not Found", http.StatusNotFound)
		return
	}

	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		http.Error(w, "Error Parsing User Data", http.StatusBadRequest)
		return
	}
}

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Error Reading Json Data", http.StatusPartialContent)
		return
	}
	token, err := auth.CreateToken(user.UserName)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	tokenCookie := &http.Cookie{
		Name:     "AuthToken",
		Value:    token,
		HttpOnly: true,
		Secure:   true,
	}
	http.SetCookie(w, tokenCookie)

	user.AuthToken = token
	user.CreatedAt = time.Now()
	newId := uuid.New()
	user.UserId = newId.String()
	err = config.CreateUser(&user)
	if err != nil {
		http.Error(w, "Error Creating in database", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

}
