package databse

import (
	"GuestList/internal/model"
	"database/sql"
	"log"
)

/* This function adds guest to a guest list table.
Arguments:
	db *sql.DB - MySQL database
	guest *model.GuestsList - guest information
Return:
	error - any error that occurred
*/
func AddGuestToList(db *sql.DB, guest *model.GuestsList) error {

	// Prepare sql query
	query, err := db.Prepare("INSERT INTO guest_list(guest_name, planned_accompanying_guests, table_id, " +
		"status, actual_accompanying_guests) VALUES ( ?, ?, ?, ?, ? )")
	if err != nil {
		log.Println(err)
		return err
	}
	defer query.Close()

	// Execute query
	_, err = query.Exec(guest.Name, guest.AccompanyingGuests, guest.TableId,
		guest.Status, -1)
	if err != nil {
		return err
	}
	log.Printf("Guest %s: successfully added to the guest list", guest.Name)
	return nil
}

/* This function checks if the table is available.
Arguments:
	db *sql.DB - database
	table int - guest information
Return:
	int - number of the available seats
	error - any error that occurred
*/
func IsTableFree(db *sql.DB, tableId int) (bool, error) {

	rows, err := db.Query("SELECT * from guest_list WHERE table_id=?", tableId)
	if err != nil {
		log.Println(err)
		return false, err
	}
	defer rows.Close()
	if rows.Next() {
		return false, nil
	}

	return true, nil
}

/* This function adds guest to a guest list table.
Arguments:
	db *sql.DB - database
	table int - guest information
Return:
	int - number of the available seats
	error - any error that occurred
*/
func GetTableCapacity(db *sql.DB, tableId int) (int, error) {
	var availableSeats int
	// Select all available seats
	rows, err := db.Query("SELECT available_seats from tables WHERE table_id=?", tableId)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	defer rows.Close()
	for rows.Next() {
		// Scan rows into variable
		if err := rows.Scan(&availableSeats); err != nil {
			log.Println(err)
			return 0, err
		}
	}
	return availableSeats, nil
}

/* This function deletes guest from the guest list table.
Arguments:
	db *sql.DB - MySQL database
	guest *model.GuestsList - guest information
Return:
	error - gives error if the guest could not be added to the guest list, or nil if the guest was added
*/
func DeleteGuestFromList(db *sql.DB, guestName string) error {
	// Prepare sql query
	query, err := db.Prepare("DELETE FROM guest_list WHERE guest_name=?")
	if err != nil {
		log.Println(err)
		return err
	}
	defer query.Close()

	// Execute query
	_, err = query.Exec(guestName)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Printf("Guest %s: successfully deleted from the guest list", guestName)
	return nil
}

/* This function gets all guest from the guest list table.
Arguments:
	db *sql.DB - MySQL database
	limit int - limit for pagination
	offset int- offset
Return:
	[]model.GuestsList - slice Guests
	error - any error that occurred
*/
func GetAllGuests(db *sql.DB, limit int, offset int) ([]model.GuestsList, error) {//([]map[string]interface{}, error) {
	var guestList []model.GuestsList
	//var guestList []map[string]interface{}
	// Select all guests
	rows, err := db.Query("SELECT guest_name, table_id, "+
		"planned_accompanying_guests from guest_list LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	guest := &model.GuestsList{}
	for rows.Next() {
		// Scan rows into Guest structure
		if err := rows.Scan(&guest.Name, &guest.TableId, &guest.AccompanyingGuests); err != nil {
			log.Println(err)
			return nil, err
		}
		// Add guest to the slice
		guestList = append(guestList, *guest)
	}
	return guestList, nil
}

/* This function gets all empty seats.
Arguments:
	db *sql.DB - MySQL database
Return:
	int - number of empty seats
	error - any error that occurred
*/
func EmptySeats(db *sql.DB) (int, error) {
	// Retrieve all arrived guests
	rows, err := db.Query("SELECT SUM(actual_accompanying_guests + 1)"+
		"FROM guest_list WHERE status=?", "ARRIVED")
	if err != nil {
		log.Println(err)
		return 0, err
	}
	defer rows.Close()

	var totalArrivedGuests int
	for rows.Next() {
		// Scan rows into variable
		if err := rows.Scan(&totalArrivedGuests); err != nil {
			log.Println(err)
			return 0, err
		}
	}
	// Retrieve the capacity of all tables
	rows, err = db.Query("SELECT SUM(available_seats) FROM tables")
	if err != nil {
		log.Println(err)
		return 0, err
	}
	defer rows.Close()

	var totalSeats int
	for rows.Next() {
		// Scan rows into variable
		if err := rows.Scan(&totalSeats); err != nil {
			log.Println(err)
			return 0, err
		}
	}

	return totalSeats - totalArrivedGuests, nil
}

/* This function gets information about invited guest.
Arguments:
	db *sql.DB - MySQL database
	guestName string - guest name
Return:
	model.GuestsList - guest information
	error - any error that occurred
*/
func GetGuestInvite(db *sql.DB, guestName string) (*model.GuestsList, error) {
	// Retrieve guest info
	rows, err := db.Query("SELECT guest_name, table_id FROM guest_list WHERE guest_name=?", guestName)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	guest := &model.GuestsList{}
	for rows.Next() {
		// Scan rows into Guest structure
		if err := rows.Scan(&guest.Name, &guest.TableId); err != nil {
			log.Println(err)
			return nil, err
		}
	}

	return guest, nil
}

/*------------------------------ Once the Party Starts ------------------------------ */

/* This function updates status of the guest to arrive.
Arguments:
	db *sql.DB - MySQL database
	guest *model.GuestsList - guest information
Return:
	error - any error that occurred
*/
func UpdateGuestStatusToArrive(db *sql.DB, guest *model.GuestsList, arrGuests int) error {
	// Let the guest in and update the status and actual arrived guests. Arrival time will get updated automatically.
	query, err := db.Prepare("UPDATE guest_list set status=?, actual_accompanying_guests=? WHERE guest_name=?")
	if err != nil {
		log.Println(err)
		return err
	}
	defer query.Close()

	// Execute query
	_, err = query.Exec("ARRIVED", arrGuests, guest.Name)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Printf("Guest %s: successfully updated from the guest list", guest.Name)
	return nil
}

/* This function gets information about the arrived guest.
Arguments:
	db *sql.DB - MySQL database
	guestName string - guest name
Return:
	*model.GuestsList - guest information
	error - any error that occurred
*/
func GetEntryFromGuestList(db *sql.DB, guestName string) (*model.GuestsList, error) {
	// Select all guests
	rows, err := db.Query("SELECT guest_name, planned_accompanying_guests, "+
		"table_id, status from guest_list WHERE guest_name=?", guestName)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	guest := &model.GuestsList{}
	for rows.Next() {
		// Scan rows into Guest structure
		if err := rows.Scan(&guest.Name, &guest.AccompanyingGuests, &guest.TableId, &guest.Status); err != nil {
			log.Println(err)
			return nil, err
		}
	}
	return guest, nil
}

/* This function gets information about all the arrived guests.
Arguments:
	db *sql.DB - MySQL database
	limit int - limit for pagination
	offset int- offset
Return:
	[]*model.GuestsList - slice of guests information
	error - any error that occurred
*/
func GetArrivedGuests(db *sql.DB, limit int, offset int) ([]model.GuestsList, error) {
	var guestList []model.GuestsList
	// Select all guests
	rows, err := db.Query("SELECT guest_name, actual_accompanying_guests, arrived_time "+
		"FROM guest_list WHERE status=? LIMIT ? OFFSET ?", "ARRIVED", limit, offset)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	guest := &model.GuestsList{}
	for rows.Next() {
		// Scan rows into Guest structure
		if err := rows.Scan(&guest.Name, &guest.AccompanyingGuests, &guest.ArrivedTime); err != nil {
			log.Println(err)
			return nil, err
		}
		// Add guest to the slice
		guestList = append(guestList, *guest)
	}

	return guestList, nil
}

