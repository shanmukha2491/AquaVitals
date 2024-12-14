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
		http.Error(w, err.Error() , http.StatusNotFound)
		return
	}

	newToken, err := auth.CreateToken(user.UserName)
	if err != nil{
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}
	user.AuthToken = newToken

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
	// Set response header
	w.Header().Set("Content-Type", "application/json")

	// Extract token from Authorization header
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		http.Error(w, "Unauthorized: No token provided", http.StatusUnauthorized)
		return
	}
	tokenString = tokenString[len("Bearer "):]

	// Decode the incoming JSON into newSensor
	var newSensor model.Sensor
	err := json.NewDecoder(r.Body).Decode(&newSensor)
	if err != nil {
		fmt.Println("Error decoding sensor data:", err)
		http.Error(w, "Invalid Sensor Data", http.StatusBadRequest)
		return
	}
	fmt.Printf("Decoded Sensor: %+v\n", newSensor)

	// Extract username from token
	username, err := auth.ExtractUsername(tokenString)
	if err != nil {
		fmt.Println("Error extracting username from token:", err)
		http.Error(w, "User not found in token", http.StatusUnauthorized)
		return
	}
	fmt.Printf("Extracted username: %s\n", username)

	// Register the sensor in the database
	err = config.RegisterSensor(newSensor, username)
	if err != nil {
		fmt.Println("Error registering sensor:", err)
		http.Error(w, "Failed to register sensor", http.StatusInternalServerError)
		return
	}
	fmt.Println("Sensor successfully registered.")

	// Send success message
	successMessage := struct {
		Success bool `json:"success"`
	}{
		Success: true,
	}

	err = json.NewEncoder(w).Encode(&successMessage)
	if err != nil {
		fmt.Println("Error encoding success message:", err)
		http.Error(w, "Error sending success message", http.StatusInternalServerError)
		return
	}
	fmt.Println("Success response sent.")
}



func FetchUser(w http.ResponseWriter, r *http.Request){
	
	w.Header().Set("Content-Type", "application/json")

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

	user, err := config.FindOneHome(username)
	if err != nil{
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	sensorData := struct{
		Sensors []model.Sensor `json:"sensors"`
	}{
		Sensors: user.Sensors,
	}
	err = json.NewEncoder(w).Encode(sensorData)
	if err != nil{
		http.Error(w, "Failed to fetch data", http.StatusNotFound)
		return
	}


}