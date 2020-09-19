package common

import (
	"GuestList/internal/databse"
	"GuestList/internal/model"
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

/*
This function adds a new guest to guest list and writes an appropriate message in response to the incoming request.
Arguments:
	resp http.ResponseWriter - HTTP response writer
	req *http.Request - HTTP request to the REST API
	db *sql.DB -  MySQL db
*/
func AddGuest(resp http.ResponseWriter, req *http.Request, db *sql.DB) {
	guest := &model.GuestsList{}
	// Get the request parameters
	params := mux.Vars(req)

	// Get the request body
	errDecoder := json.NewDecoder(req.Body).Decode(&guest)
	if errDecoder != nil {
		log.Println(errDecoder)
		encodeResponse(resp, map[string]string{"error": errDecoder.Error()}, http.StatusInternalServerError)
		return
	}

	// Retrieve name from params
	guest.Name = strings.Replace(params["name"], "+", " ", -1)
	// Set the default status for the guest
	guest.Status = "NOT_ARRIVED"

	// Check if the table is available
	free, err := databse.IsTableFree(db, *guest.TableId)
	if err!= nil {
		log.Println(errDecoder)
		encodeResponse(resp, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}
	// If the table is already reserved, return the error
	if !free {
		log.Println("The table is already reserved")
		encodeResponse(resp, map[string]string{"error": "table is already reserved"}, http.StatusBadRequest)
		return
	}

	// Get the available seats on the table
	availableSeats, err := databse.GetTableCapacity(db, *guest.TableId)
	if err != nil {
		log.Println(errDecoder)
		encodeResponse(resp, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}

	// Check if the table have enough empty seats
	if guest.AccompanyingGuests + 1 > availableSeats {
		log.Println("Not enough seats at the specified table")
		encodeResponse(resp, map[string]string{"error": "insufficient space at the specified table"}, http.StatusBadRequest)
		return
	}

	// Add the guest to a guest list
	errDB := databse.AddGuestToList(db, guest)
	// Error while adding the guest
	if errDB != nil {
		log.Println(errDB)
		encodeResponse(resp, map[string]string{"error": errDB.Error()}, http.StatusInternalServerError)
		return
	}
	// Encode the response
	encodeResponse(resp, map[string]string{"name": guest.Name}, http.StatusCreated)
}

/*
This function deletes a guest from guest list and writes an appropriate message in response to the incoming request.
Arguments:
	resp http.ResponseWriter - HTTP response writer
	req *http.Request - HTTP request to the REST API
	db *sql.DB -  MySQL db
*/
func DeleteGuest(resp http.ResponseWriter, req *http.Request, db *sql.DB) {
	// Get the request parameters
	params := mux.Vars(req)

	// Retrieve name from params. Here, the space in the name will be given as + in the REST API url.
	// Hence, we replace "+" in guest name with " ".
	guestName := strings.Replace(params["name"], "+", " ", -1)

	// Deleting guest from the guest list
	errDB := databse.DeleteGuestFromList(db, guestName)
	if errDB != nil {
		encodeResponse(resp, map[string]string{"error": errDB.Error()}, http.StatusBadRequest)
		return
	}
	// Encode the response
	encodeResponse(resp, errDB, http.StatusNoContent)
}

/*
This function gets all guest from guest list and writes an appropriate message in response to the incoming request.
Arguments:
	resp http.ResponseWriter - HTTP response writer
	req *http.Request - HTTP request to the REST API
	db *sql.DB -  MySQL db
*/
func GetGuestList(resp http.ResponseWriter, req *http.Request, db *sql.DB) {
	var limit, offset int

	// Get the request parameters
	params := req.URL.Query()
	limitVal := params.Get("limit")
	offsetVal := params.Get("offset")

	// Check if limit and offset are in the request
	if limitVal != "" {
		limit, _ = strconv.Atoi(limitVal)
	} else {
		limit = model.LIMIT
	}
	if offsetVal != "" {
		offset, _ = strconv.Atoi(offsetVal)
	} else {
		offset = model.OFFSET
	}

	// Retrieve all guests
	guestList, err := databse.GetAllGuests(db, limit, offset)
	if err != nil {
		encodeResponse(resp, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}
	// Encode the response
	encodeResponse(resp, map[string][]model.GuestsList{"guests": guestList}, http.StatusOK)
}

/*
This function gets all the guests who have arrived to the party and writes
	an appropriate message in response to the incoming request.
Arguments:
	resp http.ResponseWriter - HTTP response writer
	req *http.Request - HTTP request to the REST API
	db *sql.DB -  MySQL db
*/
func GetArrivedGuests(resp http.ResponseWriter, req *http.Request, db *sql.DB) {
	var limit, offset int

	// Get the request parameters
	params := req.URL.Query()
	limitVal := params.Get("limit")
	offsetVal := params.Get("offset")

	// Check if limit and offset is in the request
	if limitVal != "" {
		limit, _ = strconv.Atoi(limitVal)
	} else {
		limit = model.LIMIT
	}
	if offsetVal != "" {
		offset, _ = strconv.Atoi(offsetVal)
	} else {
		offset = model.OFFSET
	}

	// Retrieve arrived guests
	guestList, err := databse.GetArrivedGuests(db, limit, offset)
	if err != nil {
		encodeResponse(resp, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}
	// Encode the response
	encodeResponse(resp, map[string][]model.GuestsList{"guests": guestList}, http.StatusOK)
}

/*
This function updates a guest status to ARRIVED upon guest's arrival and writes an appropriate message
in response to the incoming request.
Arguments:
	resp http.ResponseWriter - HTTP response writer
	req *http.Request - HTTP request to the REST API
	db *sql.DB -  MySQL db
*/
func UpdateArrivedGuest(resp http.ResponseWriter, req *http.Request, db *sql.DB) {
	// Get the request parameters
	params := mux.Vars(req)

	guest := &model.GuestsList{}
	errDecoder := json.NewDecoder(req.Body).Decode(&guest)
	if errDecoder != nil {
		log.Println(errDecoder)
		encodeResponse(resp, map[string]string{"error": errDecoder.Error()}, http.StatusInternalServerError)
		return
	}
	// Retrieve name from params
	guest.Name = strings.Replace(params["name"], "+", " ", -1)
	// Get the entry from the guest list
	entry, err := databse.GetEntryFromGuestList(db, guest.Name)
	if err != nil {
		log.Println(err)
		encodeResponse(resp, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}

	// Get accompanying guests upon arrival
	arrGuests := guest.AccompanyingGuests

	// If a guest arrives with an entourage that is more than the size indicated at the guest list.
	// Check the capacity of the table and if enough seats are available allow them to come.
	if arrGuests > entry.AccompanyingGuests {
		// Get the capacity of the reserved table
		tableCapacity, err := databse.GetTableCapacity(db, *entry.TableId)
		if err != nil {
			log.Println(err)
			encodeResponse(resp, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
			return
		}

		if tableCapacity < arrGuests + 1 {
			encodeResponse(resp, map[string]string{"error": "table cannot accommodate the accompanying guests"}, http.StatusBadRequest)
			return
		}
	}

	// Update the arrival status of the guest in the guest list. This will also record the arrival time.
	errDB := databse.UpdateGuestStatusToArrive(db, guest, arrGuests)

	// Error while adding the guest
	if errDB != nil {
		encodeResponse(resp, map[string]string{"error": errDB.Error()}, http.StatusInternalServerError)
		return
	}
	// Encode the response
	encodeResponse(resp, map[string]string{"name": guest.Name}, http.StatusOK)
}

/*
This function counts all empty seats and writes an appropriate message in response to the incoming request.
Arguments:
	resp http.ResponseWriter - HTTP response writer
	req *http.Request - HTTP request to the REST API
	db *sql.DB -  MySQL db
*/
func CountEmptySeats(resp http.ResponseWriter, db *sql.DB) {

	// Get number of empty seats
	emptySeats, err := databse.EmptySeats(db)
	if err != nil {
		log.Println(err)
		encodeResponse(resp, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}
	// Encode the response
	encodeResponse(resp, map[string]int{"seats_empty": emptySeats}, http.StatusOK)
}

/*
This function generates and downloads HTML invitation.
Arguments:
	resp http.ResponseWriter - HTTP response writer
	req *http.Request - HTTP request to the REST API
	db *sql.DB -  MySQL db
*/
func GenerateInvitation (resp http.ResponseWriter, req *http.Request,  db *sql.DB) {
	// Get the request parameters
	params := mux.Vars(req)

	// Retrieve name from params
	guestName := strings.Replace(params["name"], "+", " ", -1)
	guest, err := databse.GetGuestInvite(db, guestName)
	if err != nil {
		log.Println(err)
		encodeResponse(resp, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}
	// Parse template
	tmpl, err := template.ParseFiles("templates/invitation.html")
	if err != nil {
		log.Println(err)
		encodeResponse(resp, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}
	resp.Header().Set("Content-Disposition", "attachment; filename=invitation_"+
		strings.Replace(guestName, " ", "_", -1)+".html")
	resp.Header().Set("Content-Type", req.Header.Get("Content-Type"))
	tmpl.Execute(resp, guest)
}