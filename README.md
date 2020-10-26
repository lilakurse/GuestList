# Guests List 
This project implements the guest list service (the REST API) for the year-end party!

Before the party, this service allow the functionality to add and remove guests with their entourages from the guest 
list and generate an invitation for the invited guests. 

When the party begins, guests will arrive with an entourage. This party may not be the size indicated on the guest list. 
If the guest's table can accommodate the extra people, then the whole party should be let in. Otherwise, they will be turned away.
Guests will also leave throughout the course of the party. When a guest leaves, their accompanying guests will leave with them.

At any point in the party, we should be able to know:
- Our guests at the party
- How many empty seats there are

Here, we assume that each guest (along with their entourage) will have a separate table in the party and the table will 
not be shared between two guests.

## Functionalities
The API provides the following key features:

**BEFORE PARTY**
1. Add a guest to the guest list
2. Remove a guest from the guest list
3. Get the list of guests in the guest list
4. Generate an invitation for the guest

**DURING PARTY**

5. Record the arrival of the guest to the party
6. Record guests departure from the party
7. Get a list of guests who have arrived to the party
8. Count number of empty seats at the venue

## Implementation Details
**Programming Language:** GoLang 1.14 (refer to go.mod file)

**Database:** MySQL (Version 5.7) - can be installed from [here](https://dev.mysql.com/downloads/mysql/5.7.html)

**Requirements:**
- github.com/DATA-DOG/go-sqlmock v1.5.0 (refer to go.mod file)
- github.com/go-sql-driver/mysql v1.5.0 (refer to go.mod file)
- github.com/gorilla/mux v1.8.0 (refer to go.mod file)
- github.com/stretchr/testify v1.6.1 (refer to go.mod file)
- golang-migrate - required for creating the database tables in MySQL  


## Future Improvements
In the future, I would consider the following improvements in the system:
- Run migration from the code instead of running it from the command line
- Add one more layer of business logic.
- Put all configuration in the `.env` file.
- Here, we are assuming that the tables are not shared between guests. We can modify this service to allow the table sharing
between guests and their entourage. 
- Implementing interfaces for the databases which will make the unittests easier.
- Adding more unittests and end-to-end integration tests. 


## Instructions to run the code
Please install MySQL (Version 5.7) before running the code.

1) Go to `GuestList` folder on the command prompt
2) Install `golang-migrate`
    ```
    $ brew install golang-migrate
    ```
3) Create the database tables 
    ```
     $ migrate -source file://migration -database "mysql://<user>:<pwd>@tcp(localhost:3306)/party" up
    ```
4) Build the main.go file
    ```
    $ go build main.go
    ```
5) Run the main.go file
    ```
    $ ./main
    ```

## Instructions for System Tests

**Option 1:** Go to `database` package and run the tests

1) Go to `GuestList/database/` folder

2) Run `go test`
    ```
    $ go test
    ```
**Option 2:** Run tests from `GuestList` folder
```
$ go test ./...
```

## REST API Calls

#### 1. Add a guest to the guest list
Add a given guest to the guest list

**Request URL:** http://localhost:8000/guest_list/{name}

**Request Body:** Contains table number and accompanying guests in the form of `{"table": int, "accompanying_guests": int}`

**Input Variable:** `name`: name of the guest - space is indicated using '+'

**Method:** POST

**Example:**

```
$ curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"table": 1, "accompanying_guests": 2}' \
  http://localhost:8000/guest_list/John+Smith
```

**Output:**
Returns the name of the added guest note
```
{
    "name": "John Smith"
}
```
**HTTP Response Status Code:** 201 Created

#### 2. Remove a guest from the guest list
Remove the given guest from the guest list.

**Request URL:** http://localhost:8000/guest_list/{name}

**Input Variable:** `name`: name of the guest - space is indicated using '+'

**Method:** DELETE

**Example:**

```
$ curl --header "Content-Type: application/json" \
  --request DELETE \
  http://localhost:8000/guest_list/John+Smith
```
**HTTP Response Status Code:** 204 No Content

#### 3. Get the list of guests in the guest list
Get the list of all the guests present in the guest list

**Request URL:** http://localhost:8000/guest_list/

**Method:** GET

**Example:**
```
$ curl --header "Content-Type: application/json" \
  http://localhost:8000/guest_list
```

**Output:**
Returns the guest list
```
{
    "guests": [
        {
            "name": "John Smith",
            "accompanying_guests": 2,
            "table": 1
        },
        {
            "name": "Mary Queen",
            "accompanying_guests": 3,
            "table": 2
        }
    ]
}
```
**HTTP Response Status Code:** 200 OK


#### 4. Generate an invitation for the guest
Generates an HTML file with the party invitation for the given name.

**Request URL:** http://localhost:8000/invitation/{name}

**Input Variable:** `name`: guest name

**Method:** GET

**Example:**
```
$ curl --header "Content-Type: application/json, Content-Disposition: attachment; filename=invitation_<name>.html" \
  --request GET \
  http://localhost:8000/invitation/John+Smith
```

**Output:**
Returns the HTML file.

**HTTP Response Status Code:** 200 OK

#### 5. Record the arrival of the guest to the party
Record the arrival of the guest at the party. This will also record the arrival time.

**Request URL:** http://localhost:8000/guests/{name}

**Input Variable:** `name`: guest name

**Request Body:** It contains note archived status in the form of `{"accompanying_guests": int}`

**Method:** PUT

**Example:**
```
$ curl --header "Content-Type: application/json" \
  --request PUT \
  --data '{"accompanying_guests": 2}' \
  http://localhost:8000/guests/John+Smith
```

**Output:**
Returns the update statistics from MongoDB.
```
{
    "name": "John Smith"
}
```
**HTTP Response Status Code:** 200 OK

#### 6. Record guests departure from the party
Delete the guest from the guest list upon departure.

**Request URL:** http://localhost:8000/guests/{name}

**Input Variable:** `name`: guest name

**Method:** DELETE

**Example:**
```
$ curl --header "Content-Type: application/json" \
  --request DELETE \
  http://localhost:8000/guests/John+Smith
```

**HTTP Response Status Code:** 204 No Content

#### 7. Get a list of guests who have arrived to the party
Get a list of guests who have already arrived to the party

**Request URL:** http://localhost:8000/guests/

**Method:** GET

**Example:**
```
$ curl --header "Content-Type: application/json" \
  --request GET \
  http://localhost:8000/guests/
```

**Output:**
Returns the list of guest who have arrived.
```
{
    "guests": [
        {
            "name": "John Smith",
            "accompanying_guests": 2,
            "time_arrived": "2020-09-18T16:28:44Z"
        },
        {
            "name": "Mary Queen",
            "accompanying_guests": 3,
            "time_arrived": "2020-09-18T22:32:56Z"
        }
    ]
}
```

**HTTP Response Status Code:** 200 OK

#### 8. Count number of empty seats at the venue
Count the number of empty seats at the venue 

**Request URL:** http://localhost:8000/seats_empty

**Method:** GET

**Example:**
```
$ curl --header "Content-Type: application/json" \
  --request GET \
  http://localhost:8000/seats_empty
```

**Output:**
Returns the number of empty seats.
```
{
    "seats_empty": 6
}
```
**HTTP Response Status Code:** 200 OK
