package implementation

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

const (
	// ---------- General DB constants ---------
	DB_ENGINE   = "postgres"
	DB_HOST     = "localhost"
	DB_PORT     = 5432       // standard port for Postgres
	DB_NAME     = "postgres" // we use the standard DB
	DB_USER     = "postgres"
	DB_PASSWORD = "password"
	// ---------- BIKE TABLE CONSTANTS ---------
	DB_TABLE_BIKE                      = "bike"
	DB_TABLE_BIKE_COLUMN_BIKEID        = "bikeid"
	DB_TABLE_BIKE_COLUMN_NAME          = "name"
	DB_TABLE_BIKE_COLUMN_LATITUDE      = "latitude"
	DB_TABLE_BIKE_COLUMN_LONGITUDE     = "longitude"
	DB_TABLE_BIKE_COLUMN_RESERVATIONID = "reservationid"
	// ---------- RESERVATION TABLE CONSTANTS ---------
	DB_TABLE_RESERVATION                      = "reservation"
	DB_TABLE_RESERVATION_COLUMN_RESERVATIONID = "reservationid"
	DB_TABLE_RESERVATION_COLUMN_BIKEID        = "bikeid"
	DB_TABLE_RESERVATION_COLUMN_USERNAME      = "username"
	// ---------- USER TABLE CONSTANTS ---------
	DB_TABLE_USER                 = "users"
	DB_TABLE_USER_COLUMN_USERNAME = "username"
)

/*
function to connect to Database.
returns a pointer to the connected Database
*/
func SetupDB() (*sql.DB, error) {

	// prepare connect parameters to a string
	dbinfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME,
	)
	// connect to database
	db, dbConnectError := sql.Open(DB_ENGINE, dbinfo)
	if dbConnectError != nil {
		return nil, fmt.Errorf("error while connecting to database")
	}
	// return a pointer to the connected database
	return db, nil
}

/*
returns all records from a given table without conditions
*/
func getAllRowsFromTable(db *sql.DB, tableName string) (*sql.Rows, error) {

	sqlStatement := "SELECT * FROM " + tableName

	// perform query.
	rows, dbQueryError := db.Query(sqlStatement)
	if dbQueryError != nil {
		return nil, fmt.Errorf("error retrieving all records from table " + tableName)
	}
	return rows, nil
}

/* returns all reserved bikes for a user from the reservation table */
func getBikeReservationsForUserFromDb(db *sql.DB, username string) (*sql.Rows, error) {

	// prepare Statement
	sqlStatement := getAllRecordsWithSingleConditionStatement(DB_TABLE_RESERVATION, DB_TABLE_RESERVATION_COLUMN_USERNAME)

	// Perform Query
	rows, dbQueryError := db.Query(sqlStatement, username)
	if dbQueryError != nil {
		return nil, fmt.Errorf("error retrieving bike reservations for user %v", username)
	}
	return rows, nil
}

/*
function which queries all records for a given username in the user table to verify if the user exists
returns true if the user exists in the database
*/
func userExistsInDb(db *sql.DB, username string) (bool, error) {

	queryString := getAllRecordsWithSingleConditionStatement(DB_TABLE_USER, DB_TABLE_USER_COLUMN_USERNAME)
	rows, dbQueryError := db.Query(queryString, username)

	if dbQueryError != nil {
		return false, nil
	}

	if rows.Next() {
		// record exists
		return true, nil
	}

	return false, nil
}

/*
function which queries all records for a given bikeId in the given table to verify if the bikeId exists
returns true if the bikeId exists in the table
*/
func bikeIdExistsInTable(db *sql.DB, tablename string, bikeId int) (bool, error) {
	queryString := getAllRecordsWithSingleConditionStatement(DB_TABLE_BIKE, DB_TABLE_BIKE_COLUMN_BIKEID)
	rows, dbQueryError := db.Query(queryString, bikeId)

	if dbQueryError != nil {
		return false, nil
	}

	if rows.Next() {
		// record exists
		return true, nil
	}

	return false, nil
}

/*
function, which looks up, if the bike with given bikeId is available for rent
It is available for rent, if the reservationId in the record is null
returns true, if bike is available.
*/
func bikeIsAvailableForRent(db *sql.DB, bikeId int) (bool, error) {
	queryString := getAllRecordsWithSingleConditionStatement(DB_TABLE_BIKE, DB_TABLE_BIKE_COLUMN_BIKEID)
	rows, dbQueryError := db.Query(queryString, bikeId)

	if dbQueryError != nil {
		return false, nil
	}

	if rows.Next() {
		// record exists
		// if a record exists, look for a uuid of a reservation
		tempBike := BikeImpl{}
		// fill the object
		scanError := rows.Scan(&tempBike.BikeId, &tempBike.Name, &tempBike.Latitude, &tempBike.Longitude, &tempBike.ReservationId)
		if scanError != nil {
			return false, fmt.Errorf("error scanning fields. could not scan rows of %v into bikeobject", DB_TABLE_BIKE)
		}
		// if it has a reservationId, then the bike is not available for rent
		if tempBike.ReservationId.Valid {
			return false, nil
		}

		// if it does not have a reservationId, it is available
		return true, nil
	}

	// provided bikeId does not exist
	return false, fmt.Errorf("provided bike Id does not exist")
}

/*
function, which takes a bikeId, and returns the corresponding Bike object
*/
func getBikeFromDb(db *sql.DB, bikeId int) (*BikeImpl, error) {

	queryString := getAllRecordsWithSingleConditionStatement(DB_TABLE_BIKE, DB_TABLE_BIKE_COLUMN_BIKEID)
	rows, dbQueryError := db.Query(queryString, bikeId)

	if dbQueryError != nil {
		return nil, fmt.Errorf("could not retrieve bike with %v %v from table", DB_TABLE_BIKE_COLUMN_BIKEID, bikeId, DB_TABLE_BIKE)
	}

	targetBike := BikeImpl{}
	if rows.Next() {
		// record exists
		// fill the object
		scanError := rows.Scan(&targetBike.BikeId, &targetBike.Name, &targetBike.Latitude, &targetBike.Longitude, &targetBike.ReservationId)
		if scanError != nil {
			return nil, fmt.Errorf("error scanning fields. could not scan rows of %v into bikeobject", DB_TABLE_BIKE)
		}

	}

	return &targetBike, nil
}

/*
function which creates a new record in the reservation table.

	1st param: bikeId
	2nd param: username

returns the primary key which is the newly generated uuid
*/
func createRecordInReservationTable(database *sql.DB, bikeId int, username string) (*string, error) {

	insertStatement := getInsertStmtThreeColumns(DB_TABLE_RESERVATION, DB_TABLE_RESERVATION_COLUMN_RESERVATIONID, DB_TABLE_RESERVATION_COLUMN_BIKEID, DB_TABLE_RESERVATION_COLUMN_USERNAME)

	newReservationId := uuid.New().String() // create new uuid for reservationId
	_, dbInsertError := database.Exec(insertStatement, newReservationId, bikeId, username)
	if dbInsertError != nil {
		duplicateErrorMessage := "duplicate key value"
		if strings.Contains(dbInsertError.Error(), duplicateErrorMessage) {
			return nil, fmt.Errorf("could not rent bike. User already has a rented bike")
		}

		return nil, fmt.Errorf("could not insert record into reservation Table. %v", dbInsertError)
	}

	// update bike table
	// get the update statement
	updateStmt := getUpdateStmtOneColumn(DB_TABLE_BIKE, DB_TABLE_BIKE_COLUMN_RESERVATIONID, DB_TABLE_BIKE_COLUMN_BIKEID)
	_, dbUpdateError := database.Exec(updateStmt, newReservationId, bikeId)

	if dbUpdateError != nil {
		/*
			If the record in the reservation table has been created, the update Statement needs to succeed also in the bikeTable.
			If the update does not succeed, we have an inconsistency in the database. Thus, we need to delete the newly created record which we created in the reservation table.
		*/

		// build the delete statement and delete record from database
		deleteStatement := getDeleteRowStatement(DB_TABLE_RESERVATION, DB_TABLE_RESERVATION_COLUMN_BIKEID)
		_, dbDeleteError := database.Exec(deleteStatement, bikeId)
		if dbDeleteError != nil {
			return nil, fmt.Errorf("WARNING! Inconsistency! Tried to delete newly created reservation for bike with BikeId %v but failed. Delete manually if possible. Error: %v", bikeId, dbDeleteError)
		}

		return nil, fmt.Errorf("could not insert record into reservation Table. %v", dbInsertError) // TODO: if update fails, delete entry in reservation table

	}

	return &newReservationId, nil
}

// ----------------------- Functions to get Query Strings ----------------------------------

/*
function which returns a deleteStatement for a given table and one query param for a single column
example: DELETE FROM TABLENAME WHERE CONDITION=$1
*/
func getDeleteRowStatement(dbTableName string, columnName string) string {
	deleteStatement := `delete from "` + dbTableName + `" where ` + columnName + `=$1`
	return deleteStatement
}

/*
function which returns a query Statement with one condition
example: SELECT * FROM TABLENAME WHERE CONDITION=$1
*/
func getAllRecordsWithSingleConditionStatement(dbTableName string, columnName string) string {
	queryString := `SELECT * FROM ` + dbTableName + ` WHERE ` + columnName + `=$1;`
	return queryString
}

/*
returns an insert statement string for a table with 3 columns
example: INSERT INTO TABLENAME (COLUMNNAME1, COLUMNNAME2, COLUMNNAME3) VALUES ($1, $2, $3)
*/

func getInsertStmtThreeColumns(tableName string, column1 string, column2 string, column3 string) string {
	insertInto := `insert into `
	table := `"` + tableName + `"`
	columns := `("` + column1 + `", "` + column2 + `", "` + column3 + `") `
	values := `values($1, $2, $3)`

	insertStatement := insertInto + table + columns + values

	return insertStatement
}

/*
	 returns a statement string to update a record in a table.
		1st param: table to update
		2nd param: the column to update
		3rd param: a column used to query the target row (e.g. column with primary key)
*/
func getUpdateStmtOneColumn(tableName string, columnToUpdate string, columnToQuery string) string {
	return `update "` + tableName + `" set "` + columnToUpdate + `"=$1 where "` + columnToQuery + `"=$2`
}
