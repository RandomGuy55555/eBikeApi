package implementation

import "database/sql"

/*
represents the database structure for the table "bike" in the DATABASE.
the reservationId is an uuid which can be null
Since the "Scan" method of the postgresql does not allow parsing null string values,
we use the sql.Nullstring datatype.
*/
type BikeImpl struct {
	BikeId        int            `json:"bikeid"`
	Name          string         `json:"name"`
	Latitude      string         `json:"latitude"`
	Longitude     string         `json:"longitude"`
	ReservationId sql.NullString `json:"reservationId"`
}

/*
represents the database structure for the reservation table.
the reservationId is an uuid which can be null
Since the "Scan" method of the postgresql does not allow parsing null string values,
we use the sql.Nullstring datatype.
*/
type BikeReservationImpl struct {
	ReservationId sql.NullString `json:"reservationId"`
	BikeId        int            `json:"bikeid"`
	Username      string         `json:"username"`
}
