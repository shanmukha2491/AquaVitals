package routes

import (
	"github.com/gorilla/mux"
	"github.com/shanmukha2491/AquaVitals/handlers"
)

func RegisterUserRouter(router *mux.Router) {
	router.HandleFunc("/v1/user/login", handlers.LoginHandler)
}
