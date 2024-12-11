package handlers

import (
	"encoding/json"
	"fmt"
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

	user, err = config.FindOne(user.Password, user.UserName)
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
	user.Sensors = []model.Sensor{}
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

func RegisterSensorHandler(w http.ResponseWriter, r *http.Request) {

	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		http.Error(w, "Unauthorised", http.StatusUnauthorized)
		return
	}
	tokenString = tokenString[len("Bearer "):]
	var newSensor model.Sensor
	err := json.NewDecoder(r.Body).Decode(&newSensor)
	fmt.Println(newSensor)
	if err != nil {
		fmt.Println("Error:", err)

		http.Error(w, "Invalid Sensor Data", http.StatusBadRequest)
		return
	}
	username, err := auth.ExtractUsername(tokenString)
	if err != nil {
		fmt.Println("Error:", err)

		http.Error(w, "User not found in token", http.StatusUnauthorized)
		return
	}
	err = config.RegisterSensor(newSensor, username)
	if err != nil {
		fmt.Println("Error:", err)

		http.Error(w, "Failed To register sensor", http.StatusUnauthorized)
		return
	}
}


type email struct{
	Email string `json:"email"`
}
func FetchUser(w http.ResponseWriter, r *http.Request){
	
	w.Header().Set("Content-Type", "application/json")
	var data email
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil{
		http.Error(w, "Invalid data", http.StatusBadRequest)
		return
	}
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		http.Error(w, "Unauthorised", http.StatusUnauthorized)
		return
	}
	tokenString = tokenString[len("Bearer "):]

	username, err := auth.ExtractUsername(tokenString)
	if err != nil {
		fmt.Println("Error:", err)

		http.Error(w, "User not found in token", http.StatusUnauthorized)
		return
	}

	user, err := config.FindOne(data.Email, username)
	if err != nil{
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(user)
	if err != nil{
		http.Error(w, "Failed to fetch data", http.StatusNotFound)
		return
	}


}