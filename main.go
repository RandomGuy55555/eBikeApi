// Go package
package main

/// Go fmt import
import (
	"fmt"
	"log"
	"net/http"

	"eBikeApi/services/handler"

	"github.com/gorilla/mux"
)

func main() {

	SERVERPORT := "8080"
	// Initialize router
	router := mux.NewRouter()

	// ------------------------ ENDPOINTS --------------------------------------

	// Get all available eBikes
	router.HandleFunc("/bikes/", handler.GetAllBikes).Methods("GET")

	// Get all rented eBikes from a user
	router.HandleFunc("/reservation", handler.GetBikeReservation).Queries("user", "{username}").Methods("GET")

	// Create a reservation for a bike
	router.HandleFunc("/reservation/", handler.CreateBikeReservation).Methods("POST")

	// Delete reservation for a specific bike
	router.HandleFunc("/reservation/bike/{bikeId}", handler.DeleteBikeReservation).Methods("DELETE")

	// serve the app
	fmt.Println("Listening on Localhost at %v", SERVERPORT)
	log.Fatal(http.ListenAndServe(":"+SERVERPORT, router))
}
