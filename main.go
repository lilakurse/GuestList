package main

import (
	"GuestList/config"
	"GuestList/internal/common"
	"GuestList/internal/databse"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	// Establish a connection with a DB
	db, err := databse.ConnectDB()

	if err != nil {
		log.Fatal(fmt.Sprintf("Not able to connect to DB: %v", err))
	}

	// Creates a new instance of a mux router
	router := mux.NewRouter().StrictSlash(true)

	// Set-up handlers for different requests
	// Add a guest to the guest list
	router.HandleFunc("/guest_list/{name:[a-zA-Z\\+]+}", func(w http.ResponseWriter, r *http.Request) {
		common.AddGuest(w, r, db)
	}).Methods("POST")

	// Delete a guest from the guest list
	router.HandleFunc("/guest_list/{name:[a-zA-Z\\+]+}", func(w http.ResponseWriter, r *http.Request) {
		common.DeleteGuest(w, r, db)
	}).Methods("DELETE")

	// Get the list of guests
	router.HandleFunc("/guest_list", func(w http.ResponseWriter, r *http.Request) {
		common.GetGuestList(w, r, db)
	}).Methods("GET")

	// Generate an invitation HTML file for the guest
	router.HandleFunc("/invitation/{name:[a-zA-Z\\+]+}", func(w http.ResponseWriter, r *http.Request) {
		common.GenerateInvitation(w, r, db)}).Methods("GET")

	// Update the status of the guest upon arrival
	router.HandleFunc("/guests/{name:[a-zA-Z\\+]+}", func(w http.ResponseWriter, r *http.Request) {
		common.UpdateArrivedGuest(w, r, db)
	}).Methods("PUT")

	// Delete the guest upon departure
	router.HandleFunc("/guests/{name:[a-zA-Z\\+]+}", func(w http.ResponseWriter, r *http.Request) {
		common.DeleteGuest(w, r, db)
	}).Methods("DELETE")

	// List guests which have arrived at the party
	router.HandleFunc("/guests", func(w http.ResponseWriter, r *http.Request) {
		common.GetArrivedGuests(w, r, db)
	}).Methods("GET")

	// Get the number of empty seats
	router.HandleFunc("/seats_empty", func(w http.ResponseWriter, r *http.Request) {
		common.CountEmptySeats(w, db)
	}).Methods("GET")

	log.Fatal(http.ListenAndServe(config.API_PORT, router))
}
