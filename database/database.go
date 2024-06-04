package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type Databaser interface {
	InitDatabase()
}

type database struct{}

func NewDatabase() database {
	return database{}
}

func InitDatabase() {
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
