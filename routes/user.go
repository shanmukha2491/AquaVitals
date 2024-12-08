package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/shanmukha2491/AquaVitals/config"
	"github.com/shanmukha2491/AquaVitals/handlers"
	auth "github.com/shanmukha2491/AquaVitals/middlewares"
)

func RegisterUserRouter(router *mux.Router) {
	config.UserCollection(config.Client)
	
	router.Handle("/v1/user/login",
		auth.AuthorizationMiddleware(http.HandlerFunc(handlers.LoginHandler))).Methods(http.MethodPost)

	router.HandleFunc("/v1/user/create",
		handlers.SignUpHandler).Methods(http.MethodPost)

	router.Handle("/v1/user/register_sensor",
		auth.AuthorizationMiddleware(http.HandlerFunc(handlers.RegisterSensorHandler))).Methods(http.MethodPut)
}
