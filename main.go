package main

import (
	"log"
	"net/http"
	"os"
	c "uts/controller"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/rooms", c.GetAllRooms).Methods("GET")
	router.HandleFunc("/roomDetails", c.GetDetailRoom).Methods("GET")

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowCredentials: true,
	})

	http.Handle("/", router)

	log.Println("Starting " + os.Getenv("APP_NAME"))
	log.Println("Connected to port 8080")
	log.Fatal(http.ListenAndServe(":8080", corsHandler.Handler(router)))
}
