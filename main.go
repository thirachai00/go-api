package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"gotest/model"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func setJsonResp(message []byte, httpCode int, res http.ResponseWriter) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(httpCode)
	res.Write(message)
}

func customer(res http.ResponseWriter, req *http.Request) {

	db, _ := sql.Open("sqlite3", "./customers.db") // Open the created SQLite File
	defer db.Close()

	if req.Method == "GET" { //Get all
		var customers []model.Customer
		row, err := db.Query("SELECT * FROM Customers ORDER BY id")
		if err != nil {
			message := []byte(`{"message": "query error"}`)
			setJsonResp(message, http.StatusInternalServerError, res)
			return
		}
		defer row.Close()

		for row.Next() { // Iterate and fetch the records from result cursor
			var id int
			var name string
			var age int
			row.Scan(&id, &name, &age)
			log.Println("Customer: ", name, " ", age)
			customers = append(customers, model.Customer{ID: id, Name: name, Age: age})
		}

		response := model.JsonResponse{
			Message: "SUCCESS",
			Data:    customers,
		}

		customerJson, err := json.Marshal(response)
		if err != nil {
			message := []byte(`{"message": "Error marshalling data"}`)
			setJsonResp(message, http.StatusInternalServerError, res)
			return
		}

		setJsonResp(customerJson, http.StatusOK, res)

	} else { //Create
		var customer model.Customer
		payload := req.Body

		defer req.Body.Close()

		err := json.NewDecoder(payload).Decode(&customer)
		if err != nil {
			message := []byte(`{"message": "Error marshalling data"}`)
			setJsonResp(message, http.StatusInternalServerError, res)
			return
		}

		insertCustomerSQL := `INSERT INTO Customers(name, age) VALUES (?, ?)`
		statement, err := db.Prepare(insertCustomerSQL) // Prepare statement.
		// This is good to avoid SQL injections
		if err != nil {
			message := []byte(`{"message": "Error prepare database"}`)
			setJsonResp(message, http.StatusInternalServerError, res)
			return
		}
		_, err = statement.Exec(customer.Name, customer.Age)
		if err != nil {
			message := []byte(`{"message": "insert database error"}`)
			setJsonResp(message, http.StatusInternalServerError, res)
			return
		}

		message := []byte(`{"message": "SUCCESS"}`)

		setJsonResp(message, http.StatusCreated, res)
	}
}

func customerById(res http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" && req.Method != "DELETE" && req.Method != "PUT" {
		message := []byte(`{"message": "Invalid HTTP Method"}`)
		setJsonResp(message, http.StatusMethodNotAllowed, res)
		return
	}

	if _, ok := req.URL.Query()["id"]; !ok {
		message := []byte(`{"message": "Please provide id"}`)
		setJsonResp(message, http.StatusBadRequest, res)
		return
	}

	id := req.URL.Query().Get("id")

	db, _ := sql.Open("sqlite3", "./customers.db") // Open the created SQLite File
	defer db.Close()

	if req.Method == "GET" {
		row, err := db.Query("SELECT * FROM Customers WHERE id = ? ", id)
		if err != nil {
			fmt.Printf("error : %v\n", err)
			message := []byte(`{"message": "Error prepare database"}`)
			setJsonResp(message, http.StatusInternalServerError, res)
			return
		}
		defer row.Close()

		customer := model.Customer{}
		for row.Next() { // Iterate and fetch the records from result cursor
			var id int
			var name string
			var age int
			row.Scan(&id, &name, &age)
			customer = model.Customer{ID: id, Name: name, Age: age}
		}

		response := model.JsonResponse{
			Message: "SUCCESS",
			Data:    customer,
		}

		customerJson, err := json.Marshal(response)

		if err != nil {
			message := []byte(`{"message": "Error marshalling data"}`)
			setJsonResp(message, http.StatusInternalServerError, res)
			return
		}

		setJsonResp(customerJson, http.StatusOK, res)

	}

	if req.Method == "DELETE" {
		insertCustomerSQL := `DELETE FROM Customers WHERE id = ?`
		statement, err := db.Prepare(insertCustomerSQL) // Prepare statement.
		if err != nil {
			fmt.Printf("error : %v\n", err)
			message := []byte(`{"message": "Error prepare database"}`)
			setJsonResp(message, http.StatusInternalServerError, res)
			return
		}
		_, err = statement.Exec(id)
		if err != nil {
			fmt.Printf("error : %v\n", err)
			message := []byte(`{"message": "Error execute database"}`)
			setJsonResp(message, http.StatusInternalServerError, res)
			return
		}

		message := []byte(`{"message": "SUCCESS"}`)
		setJsonResp(message, http.StatusOK, res)
		return

	}

	if req.Method == "PUT" {
		payload := req.Body

		defer req.Body.Close()

		var updateCustomer model.Customer

		err := json.NewDecoder(payload).Decode(&updateCustomer)
		if err != nil {
			message := []byte(`{"message": "Error marshalling data"}`)
			setJsonResp(message, http.StatusInternalServerError, res)
			return
		}

		fmt.Printf("name : %v | age : %v\n", updateCustomer.Name, updateCustomer.Age)

		customerSQL := `UPDATE Customers SET name = ? , age = ? WHERE id = ?`
		statement, err := db.Prepare(customerSQL) // Prepare statement.
		if err != nil {
			fmt.Printf("error : %v\n", err)
			message := []byte(`{"message": "Error prepare database"}`)
			setJsonResp(message, http.StatusInternalServerError, res)
			return
		}

		_, err = statement.Exec(updateCustomer.Name, updateCustomer.Age, id)
		if err != nil {
			fmt.Printf("error : %v\n", err)
			message := []byte(`{"message": "Error execute database"}`)
			setJsonResp(message, http.StatusInternalServerError, res)
			return
		}

		message := []byte(`{"message": "SUCCESS"}`)
		setJsonResp(message, http.StatusOK, res)
		return
	}

}

func main() {
	initDatabase()
	http.HandleFunc("/", func(res http.ResponseWriter, r *http.Request) {
		message := []byte(`{"message": "Server up and running"}`)
		setJsonResp(message, http.StatusOK, res)
	})

	http.HandleFunc("/customers", customer)
	http.HandleFunc("/customer", customerById)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initDatabase() {
	os.Remove("customers.db") // I delete the file to avoid duplicated records.
	// SQLite is a file based database.

	log.Println("Creating customers.db...")
	file, err := os.Create("customers.db") // Create SQLite file
	if err != nil {
		log.Fatal(err.Error())
	}
	file.Close()
	log.Println("customers.db created")

	sqliteDatabase, _ := sql.Open("sqlite3", "./customers.db") // Open the created SQLite File
	defer sqliteDatabase.Close()                               // Defer Closing the database
	createTable(sqliteDatabase)                                // Create Database Tables

	// INSERT RECORDS
	insertStudent(sqliteDatabase, "Liana Kim", 20)
	insertStudent(sqliteDatabase, "Glen Rangel", 21)
	insertStudent(sqliteDatabase, "Martin Martins", 22)
	insertStudent(sqliteDatabase, "Alayna Armitage", 23)
	insertStudent(sqliteDatabase, "Marni Benson", 24)
	insertStudent(sqliteDatabase, "Derrick Griffiths", 25)
	insertStudent(sqliteDatabase, "Leigh Daly", 26)
	insertStudent(sqliteDatabase, "Marni Benson", 27)
	insertStudent(sqliteDatabase, "Klay Correa", 28)

	// DISPLAY INSERTED RECORDS
	displayCustomers(sqliteDatabase)
}

func createTable(db *sql.DB) {
	createCustomerTableSQL := `CREATE TABLE Customers (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		"name" TEXT,
		"age" integer
	  );` // SQL Statement for Create Table

	log.Println("Create Customers table...")
	statement, err := db.Prepare(createCustomerTableSQL) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec() // Execute SQL Statements
	log.Println("Customers table created")
}

// We are passing db reference connection from main to our method with other parameters
func insertStudent(db *sql.DB, name string, age int) {
	log.Println("Inserting customer record ...")
	insertCustomerSQL := `INSERT INTO Customers(name, age) VALUES (?, ?)`
	statement, err := db.Prepare(insertCustomerSQL) // Prepare statement.
	// This is good to avoid SQL injections
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(name, age)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func displayCustomers(db *sql.DB) {
	row, err := db.Query("SELECT * FROM Customers ORDER BY name")
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	for row.Next() { // Iterate and fetch the records from result cursor
		var id int
		var name string
		var age int
		row.Scan(&id, &name, &age)
		log.Println("Customer: ", name, " ", age)
	}
}
