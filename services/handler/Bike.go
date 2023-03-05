package handler

import (
	"eBikeApi/services/implementation"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// returns all available bikes from the database
func GetAllBikes(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Getting all eBikes from the database")

	allBikes, getAllBikesError := implementation.GetAllBikes()
	if getAllBikesError != nil {
		getAllBikesErrMsg := fmt.Errorf("could not retrieve all bikes. %v", getAllBikesError)
		JSONError(w, getAllBikesErrMsg, http.StatusInternalServerError)
		return
	}

	// transform the data to Response Object
	getAllBikesResponse := transformBikeImplToGetBikeResponse(allBikes)

	//return  response
	json.NewEncoder(w).Encode(getAllBikesResponse)
}

// handler method to get all bike reservations for a specific user
func GetBikeReservation(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Getting Bike Reservation for a specific user")

	vars := mux.Vars(r)
	//username := vars["username"] // TODO: retrieve username through keycloak
	username := vars["username"]

	// if username is not provided throw error
	if username == "" {
		bikeIdMissingMsg := fmt.Errorf("mandatory username not provided")
		JSONError(w, bikeIdMissingMsg, http.StatusBadRequest)
		return
	}

	// call GetBikeReservation implementation
	bikeReservations, getBikeReservationError := implementation.GetBikeReservation(username)
	if getBikeReservationError != nil {
		getBikeReservationErrMsg := fmt.Errorf("could not get bike reservation. %v", getBikeReservationError)
		JSONError(w, getBikeReservationErrMsg, http.StatusInternalServerError)
		return
	}
	// todo: improve
	tempArray := []implementation.BikeImpl{}

	tempArray = append(tempArray, *bikeReservations)

	bikeReservationResponse := transformBikeImplToGetBikeResponse(&tempArray)

	json.NewEncoder(w).Encode(bikeReservationResponse)
}

/*
	 handler method to create a reservation for a bike
		takes a http body with following values
		"username" : username
		"bikeId"
		since an id provider like keycloak does not exist we can not deliver a token.
		so unfortunately we need to provide the username in the body
*/
func CreateBikeReservation(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Creating bike reservation")

	// init variable with a struct
	var bikeReservationRequest implementation.BikeReservationImpl

	// read the request body and parse it into the struct
	readRequestError := ReadRequestBody(r.Body, &bikeReservationRequest)
	if readRequestError != nil {
		readRequestErrorMsg := fmt.Errorf("error while reading request. %v", readRequestError)
		JSONError(w, readRequestErrorMsg, http.StatusInternalServerError)
		return
	}

	// call the implementation to reserve a bike (create bike reservation)
	reserveBikeResponse, reserveBikeError := implementation.ReserveBike(bikeReservationRequest)
	if reserveBikeError != nil {
		reserveBikeErrMsg := fmt.Errorf("could not create bike reservation. %v", reserveBikeError)
		JSONError(w, reserveBikeErrMsg, http.StatusInternalServerError)
		return
	}

	JsonSuccessResponse(w, *reserveBikeResponse)
}

/*
	 handler method to delete a bike reservation
		parmameters required:
		- bikeId
*/
func DeleteBikeReservation(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Deleting bike reservation")

	vars := mux.Vars(r)
	//username := vars["username"] // TODO: retrieve username through keycloak
	bikeIdAsString := vars["bikeId"]

	// if the bikeId is not provided, throw error
	if bikeIdAsString == "" {
		bikeIdMissingMsg := fmt.Errorf("mandatory bikeId not provided")
		JSONError(w, bikeIdMissingMsg, http.StatusBadRequest)
		return
	}

	// parse the bikeId string to an Integer
	bikeId, parseErr := strconv.Atoi(bikeIdAsString)
	if parseErr != nil {
		stringToIntParseErr := fmt.Errorf("error parsing string to int. %v", parseErr)
		JSONError(w, stringToIntParseErr, http.StatusInternalServerError)
		return
	}

	// call implementation method to delete a bike reservation
	deleteBikeReservationError := implementation.DeleteBikeReservation(bikeId)
	if deleteBikeReservationError != nil {
		deleteBikeReservationErrMsg := fmt.Errorf("could not return bike. %v", deleteBikeReservationError)
		JSONError(w, deleteBikeReservationErrMsg, http.StatusInternalServerError)
		return
	}

	// return message that the deletion was successful
	JsonSuccessResponse(w, "Successfully deleted bike reservation")
}
