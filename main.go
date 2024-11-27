package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/shanmukha2491/AquaVitals/routes"
)

func main() {
	router := mux.NewRouter()

	// Register application routes
	routes.RegisterUserRouter(router)

	// Start the server
	log.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))

	// token, _ := auth.CreateToken("Shanmukha")
	// fmt.Println("Token:", token)
}
