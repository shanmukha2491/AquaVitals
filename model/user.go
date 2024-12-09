package model

import "time"

type Sensor struct {
	SensorId    string            `json:"sensor_id"`
	SensorName  string            `json:"name"`
	Location    string            `json:"location"`
	Data        map[string]string `json:"data"`
	DataBaseUrl string            `json:"database_url"`
	WebApiToken string            `json:"web_token"`
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
