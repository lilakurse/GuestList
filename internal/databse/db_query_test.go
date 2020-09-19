package databse

import (
	"GuestList/internal/model"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Test adding a guest to the guest list
func TestAddGuestToList(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	tableID := 1
	guest := &model.GuestsList{
		Name: "John Smith",
		TableId: &tableID,
		AccompanyingGuests: 2,
		Status: "NOT_ARRIVED",
	}
	prep := mock.ExpectPrepare("^INSERT INTO guest_list*")
	prep.ExpectExec().
		WithArgs("John Smith", 2, &tableID, "NOT_ARRIVED", -1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	defer db.Close()

	err = AddGuestToList(db, guest)

	assert.Equal(t, nil, err, "Expected no error")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}

}

// Test getting the seating capacity of a table
func TestGetTableCapacity(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	tableID := 1
	rows := sqlmock.NewRows([]string{"available_seats"}).AddRow(9)

	//mock.ExpectQuery("^SELECT (.+) FROM menu_link_content_data*")
	mock.ExpectQuery(
		`^SELECT available_seats from tables*`).
		WithArgs(tableID).WillReturnRows(rows)
	availableSeats, _ := GetTableCapacity(db, 1)
	assert.Equal(t, 9, availableSeats,"Expected different number of table capacity")
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}

}

// Test getting all the guests
func TestGetAllGuests(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Here we are creating rows in our mocked database.
	rows := sqlmock.NewRows([]string{"guest_name", "table_id", "planned_accompanying_guests"}).
		AddRow("John Smith", 1, 2).
		AddRow("Brad Pitt", 2, 4)

	mock.ExpectQuery(
		`^SELECT guest_name, table_id, planned_accompanying_guests from guest_list*`).
		WithArgs(10, 0).WillReturnRows(rows)
	guestList, _ := GetAllGuests(db, 10, 0)

	assert.Equal(t, 2, len(guestList),"Expected different number of guests")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}

}

// Test deleting a guest from the guest list
func TestDeleteGuestFromList(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	prep := mock.ExpectPrepare("^DELETE FROM guest_list WHERE guest_name*")
	prep.ExpectExec().
		WithArgs("John Smith").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = DeleteGuestFromList(db, "John Smith")

	assert.Equal(t, nil, err, "Expected no error")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}

}

// Test updating a guest status upon arrival
func TestUpdateGuestStatusToArrive(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	tableId := 1
	guest := &model.GuestsList{
	 	Name: "Vanessa Smith",
	 	TableId: &tableId,
	 	AccompanyingGuests: 4,
	 }
	arrivingAccompanyingGuests := 5
	prep := mock.ExpectPrepare("^UPDATE guest_list*")
	prep.ExpectExec().
		WithArgs("ARRIVED", arrivingAccompanyingGuests, guest.Name).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = UpdateGuestStatusToArrive(db, guest, arrivingAccompanyingGuests)

	assert.Equal(t, nil, err, "Expected no error")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}

}
