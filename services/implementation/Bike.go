package implementation

import (
	"fmt"
)

/*
Implementation method to retrieve all bikes from the Database
*/
func GetAllBikes() (*[]BikeImpl, error) {
	// connect to database
	db, dbConnectError := SetupDB()
	if dbConnectError != nil {
		return nil, dbConnectError
	}
	defer db.Close() // close connection to DB after finishing method

	// Get all bikes from the database
	rows, getAllRowsFromTableErr := getAllRowsFromTable(db, DB_TABLE_BIKE)
	if getAllRowsFromTableErr != nil {
		return nil, getAllRowsFromTableErr
	}

	var arrayOfBikes []BikeImpl
	// For each record...
	for rows.Next() {
		// create a new Bike Object
		tempBike := BikeImpl{}
		// fill the object
		// reservationID is a nullstring type and will be converted to "rented" (boolean) in the transform method
		scanError := rows.Scan(&tempBike.BikeId, &tempBike.Name, &tempBike.Latitude, &tempBike.Longitude, &tempBike.ReservationId)

		if scanError != nil {
			return nil, fmt.Errorf("error scanning fields. could not scan rows of %v into bikeobject", DB_TABLE_BIKE)
		}

		arrayOfBikes = append(arrayOfBikes, tempBike)
	}

	return &arrayOfBikes, nil
}

/* Implementation method to Get the bike reservation from a specific user */
func GetBikeReservation(username string) (*BikeImpl, error) {
	// connect to database
	db, dbConnectError := SetupDB()
	if dbConnectError != nil {
		return nil, dbConnectError
	}
	defer db.Close() // close connection to DB after finishing method

	// Get all bikes from the database
	reservationRecords, getAllRowsFromTableErr := getBikeReservationsForUserFromDb(db, username)
	if getAllRowsFromTableErr != nil {
		return nil, getAllRowsFromTableErr
	}

	var arrayOfBikeReservations []BikeReservationImpl
	// For each record...
	for reservationRecords.Next() {
		// create a new Reservation Object
		tempReservation := BikeReservationImpl{}
		// fill the object
		scanError := reservationRecords.Scan(&tempReservation.ReservationId, &tempReservation.BikeId, &tempReservation.Username)

		if scanError != nil {
			return nil, fmt.Errorf("error scanning fields. could not scan rows of %v into BikeReservation Object", DB_TABLE_RESERVATION)
		}

		arrayOfBikeReservations = append(arrayOfBikeReservations, tempReservation)
	}

	// if no reservations are found in the reservationTable, return empty array
	if len(arrayOfBikeReservations) == 0 {
		return &BikeImpl{}, nil
	}

	// if there is more than one reservation, return error because this means inconsistency, since a user is only allowed to have one bike.
	if len(arrayOfBikeReservations) > 1 {
		return nil, fmt.Errorf("found more than one reservation for user %v. Data Inconsistency", username)
	}

	// there is just one record... Retrieve the record from the bike table to retrieve all information of the bike
	targetReservation := arrayOfBikeReservations[0]

	// we got the targetReservation record from the reservation table. Now retrieve all bike information from the bike table via the bikeId
	targetBike, getBikeFromDbError := getBikeFromDb(db, targetReservation.BikeId)
	if getBikeFromDbError != nil {
		return nil, getBikeFromDbError
	}

	return targetBike, nil
}

/*
	 This method is for creating a bike reservation in the reservation table
		It creates a record in the reservation table.
		After successful record creation, it updates the column for the reserved bike in the bike table.
		If updating the bike record fails, the record in the reservation table will be deleted. otherwise we have an inconsistency.
*/
func ReserveBike(bikeReservationRequest BikeReservationImpl) (*string, error) {

	username := bikeReservationRequest.Username
	bikeId := bikeReservationRequest.BikeId

	// if the username is missing, throw error
	if username == "" {
		usernameMissingError := fmt.Errorf("no username provided. Bike reservation process failed")
		return nil, usernameMissingError
	}

	// if bikeId == "" { // TODO: How to check that it is not provided
	// 	bikeIdMissingError := fmt.Errorf("no bikeId provided. Bike reservation process failed")
	// 	return nil, bikeIdMissingError
	// }

	// connect to database
	db, dbConnectError := SetupDB()
	if dbConnectError != nil {
		return nil, dbConnectError
	}
	defer db.Close() // close connection to DB after finishing method

	// verify if user exists in the database
	userRecordExists, userExistsInDbError := userExistsInDb(db, username)
	if userExistsInDbError != nil {
		return nil, userExistsInDbError
	}

	if !userRecordExists {
		return nil, fmt.Errorf("provided username does not exist in database")
	}

	// verify if provided bikeId exists in the database
	bikeIdExistsInBikeTable, bikeIdExistsInDbError := bikeIdExistsInTable(db, DB_TABLE_BIKE, bikeId)
	if bikeIdExistsInDbError != nil {
		return nil, bikeIdExistsInDbError
	}

	if !bikeIdExistsInBikeTable {
		return nil, fmt.Errorf("provided bikeId does not exist in database")
	}

	// verify if provided bikeId is available for rent
	bikeIsAvailable, bikeIdExistsInDbError := bikeIsAvailableForRent(db, bikeId)
	if bikeIdExistsInDbError != nil {
		return nil, bikeIdExistsInDbError
	}

	if !bikeIsAvailable {
		return nil, fmt.Errorf("provided bikeId is not available for rent")
	}

	//create reservation by inserting it into reservation table
	createdReservationId, createReservationRecordErr := createRecordInReservationTable(db, bikeId, username)
	if createReservationRecordErr != nil {
		return nil, fmt.Errorf("could not insert record into reservation Table. %v", createReservationRecordErr)
	}

	return createdReservationId, nil
}

// deletes a Bike reservation in the reservation table for given bikeId
// there is no need to update the bike table, since database is set to "ON DELETE SET NULL"
func DeleteBikeReservation(bikeId int) error {

	// connect to DB
	db, dbConnectError := SetupDB()
	if dbConnectError != nil {
		return dbConnectError
	}
	defer db.Close() // close connection to DB after finishing method.

	// verify if provided bikeId exists in the bike table
	bikeIdExistsInBikeTable, bikeIdExistsInDbError := bikeIdExistsInTable(db, DB_TABLE_BIKE, bikeId)
	if bikeIdExistsInDbError != nil {
		return bikeIdExistsInDbError
	}

	if !bikeIdExistsInBikeTable {
		return fmt.Errorf("provided bikeId does not exist in database")
	}

	// verify if provided bikeId is available for rent
	bikeIsAvailable, bikeIsAvailableForRentError := bikeIsAvailableForRent(db, bikeId)
	if bikeIdExistsInDbError != nil {
		return bikeIsAvailableForRentError
	}

	// if bike is available, there is no reservation to delete
	if bikeIsAvailable {
		return fmt.Errorf("provided bikeId is not renteds so there is no reservation to delete")
	}

	// build the delete statement and delete record from database
	deleteStatement := getDeleteRowStatement(DB_TABLE_RESERVATION, DB_TABLE_RESERVATION_COLUMN_BIKEID)
	_, dbDeleteError := db.Exec(deleteStatement, bikeId)
	if dbDeleteError != nil {
		return fmt.Errorf("could not delete record into reservation Table. %v", dbDeleteError)
	}

	return nil
}
