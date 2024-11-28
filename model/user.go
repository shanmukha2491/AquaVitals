package model

import "time"

type Sensor struct {
	SensorName string
	Location   time.Location
	Data       map[string]string
}
type User struct {
	UserId    string    `json:"user_id"`
	UserName  string    `json:"username"`
	Password  string    `json:"password"`
	Email     string    `json:"email"`
	AuthToken string    `json:"authToken"`
	Sensors   []Sensor  `json:"sensor_data"`
	CreatedAt time.Time `json:"created_at"`
}
