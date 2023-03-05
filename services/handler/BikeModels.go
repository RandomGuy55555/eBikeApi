package handler

import (
	"eBikeApi/services/implementation"
)

/* struct used to return data from bikes in the DB as JSON response
 */
type GetBikesResponse struct {
	BikeId    int    `json:"bikeId"`
	Name      string `json:"name"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
	Rented    bool   `json:"rented"`
}

/*
since there is a difference between the data in the database, and the data we want to show in the UI,
we distinguish between them.
This function transforms the struct from the implementation layer to the struct for the JSON Response
*/
func transformBikeImplToGetBikeResponse(bikeArray *[]implementation.BikeImpl) *[]GetBikesResponse {

	var getBikeResponse []GetBikesResponse

	for _, bike := range *bikeArray {

		// null strings are not supported by the Scan method, so we need to calculate the boolean
		rented := false
		// if it is a valid reservationId, set rented to true
		if bike.ReservationId.Valid {
			rented = true
		}

		tempBike := GetBikesResponse{
			BikeId:    bike.BikeId,
			Name:      bike.Name,
			Latitude:  bike.Latitude,
			Longitude: bike.Longitude,
			Rented:    rented,
		}

		getBikeResponse = append(getBikeResponse, tempBike)
	}
	return &getBikeResponse
}

func transformBikeImplObjectToGetBikeResponse(bikeImplObject *implementation.BikeImpl) GetBikesResponse {

	var getBikeResponse GetBikesResponse

	if bikeImplObject != nil {
		// only process object if it is valid
		if isValidBikeObject(*bikeImplObject) {
			// null strings are not supported by the Scan method, so we need to calculate the boolean
			rented := false
			// if it is a valid reservationId, set rented to true
			if bikeImplObject.ReservationId.Valid {
				rented = true
			}

			tempBike := GetBikesResponse{
				BikeId:    bikeImplObject.BikeId,
				Name:      bikeImplObject.Name,
				Latitude:  bikeImplObject.Latitude,
				Longitude: bikeImplObject.Longitude,
				Rented:    rented,
			}

			getBikeResponse = tempBike
		}
	}

	return getBikeResponse
}

/*
this function checks if an object from the database is valid by checking if the string fields are not empty
this function prevents that objects are returned with an integer value of 0 and boolean value false
*/
func isValidBikeObject(bikeObject implementation.BikeImpl) bool {
	if bikeObject.Name == "" && bikeObject.Latitude == "" && bikeObject.Longitude == "" {
		return false
	}
	return true
}
