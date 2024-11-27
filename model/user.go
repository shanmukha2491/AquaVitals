package model

import "time"

type Sensor struct {
	SensorName string
	Location   time.Location
	Data       map[string]string
}
type User struct {
	UserId    string    `json:"id"`
	UserName  string    `json:"username"`
	Password  string    `json:"password"`
	Email     string    `json:"email"`
	AuthToken string    `json:"authToken"`
	Sensors   []Sensor  `json:"sensor_data"`
	CreateAt  time.Time `json:"create_at"`
}
