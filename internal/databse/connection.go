package databse

import (
	"GuestList/config"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

/* This function takes care of connecting to the given Mysql and returns the required database
Returns:
	*sql.DB - database
	error - HTTP status
*/
func ConnectDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", config.MYSQL_DSN+config.MYSQL_DATABASE+"?parseTime=true")
	if err != nil {
		log.Fatal(err.Error())
	}
	// Check if the Mysql is accessible
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	log.Println("Connected to MySQL database")

	return db, nil
}
