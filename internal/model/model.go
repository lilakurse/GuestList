package model

import "time"

const (
	LIMIT  = 100
	OFFSET = 0
)

// Model for Guests List
type GuestsList struct {
	Name               string    `json:"name"`  					// Guest name
	AccompanyingGuests int       `json:"accompanying_guests"`		// Number of accompanying guests
	TableId            *int       `json:"table,omitempty"`			// Table ID
	Status             string    `json:"-"`							// ARRIVED/NOT_ARRIVED
	ArrivedTime        *time.Time `json:"time_arrived,omitempty"`	// time of arrival in the party
}
