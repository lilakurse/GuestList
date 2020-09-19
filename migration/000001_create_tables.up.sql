CREATE DATABASE IF NOT EXISTS party;

SET @@time_zone = 'SYSTEM';

CREATE TABLE IF NOT EXISTS tables(
   table_id serial,
   available_seats INT NOT NULL,
   PRIMARY KEY(table_id)
);

CREATE TABLE IF NOT EXISTS guest_list(
   guest_id serial,
   guest_name VARCHAR (50) UNIQUE NOT NULL,
   planned_accompanying_guests INT NOT NULL,
   table_id BIGINT UNSIGNED,
   status VARCHAR(20) NOT NULL,
   actual_accompanying_guests INT NOT NULL,
   arrived_time DATETIME ON UPDATE CURRENT_TIMESTAMP,
   PRIMARY KEY (guest_id),
   FOREIGN KEY (table_id) REFERENCES tables(table_id)
);