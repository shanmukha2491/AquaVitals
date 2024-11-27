package handlers

import (
	"fmt"
	"net/http"

	auth "github.com/shanmukha2491/AquaVitals/middlewares"
)

type User struct{}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		http.Error(w, "Unauthorised", http.StatusUnauthorized)
		return
	}
	tokenString = tokenString[len("Bearer "):]
	err := auth.VerifyToken(tokenString)
	if err != nil {
		http.Error(w, "Wrong Auth Token", http.StatusUnauthorized)
		return
	}

	fmt.Fprint(w, "Welcome to the the protected area")

}
