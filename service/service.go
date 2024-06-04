package service

import (
	"database/sql"
	"gotest/model"
)

type Servicer interface {
	GetAllCustomer(db *sql.DB) ([]model.Customer, error)
	CreateCustomer(req model.Customer, db *sql.DB) error
	GetCustomerById(id string, db *sql.DB) (model.Customer, error)
	DeleteCustomer(id string, db *sql.DB) error
	UpdateCustomer(id string, req model.Customer, db *sql.DB) error
}

type service struct {
}

func NewService() service {
	return service{}
}

func GetAllCustomer(db *sql.DB) ([]model.Customer, error) {
	var customers []model.Customer
	row, err := db.Query("SELECT * FROM Customers ORDER BY id")
	if err != nil {
		return customers, err
	}
	defer row.Close()

	for row.Next() { // Iterate and fetch the records from result cursor
		var id int
		var name string
		var age int
		row.Scan(&id, &name, &age)
		customers = append(customers, model.Customer{ID: id, Name: name, Age: age})
	}

	return customers, nil
}

func CreateCustomer(req model.Customer, db *sql.DB) error {
	insertCustomerSQL := `INSERT INTO Customers(name, age) VALUES (?, ?)`
	statement, err := db.Prepare(insertCustomerSQL) // Prepare statement.
	// This is good to avoid SQL injections
	if err != nil {
		return err
	}
	_, err = statement.Exec(req.Name, req.Age)
	if err != nil {
		return err
	}

	return nil
}

func GetCustomerById(id string, db *sql.DB) (model.Customer, error) {
	customer := model.Customer{}
	row, err := db.Query("SELECT * FROM Customers WHERE id = ? ", id)
	if err != nil {
		return customer, err
	}
	defer row.Close()

	for row.Next() { // Iterate and fetch the records from result cursor
		var id int
		var name string
		var age int
		row.Scan(&id, &name, &age)
		customer = model.Customer{ID: id, Name: name, Age: age}
	}

	return customer, nil
}

func DeleteCustomer(id string, db *sql.DB) error {
	insertCustomerSQL := `DELETE FROM Customers WHERE id = ?`
	statement, err := db.Prepare(insertCustomerSQL) // Prepare statement.
	if err != nil {
		return err
	}
	_, err = statement.Exec(id)
	if err != nil {
		return err
	}

	return nil
}

func UpdateCustomer(id string, req model.Customer, db *sql.DB) error {
	customerSQL := `UPDATE Customers SET name = ? , age = ? WHERE id = ?`
	statement, err := db.Prepare(customerSQL) // Prepare statement.
	if err != nil {
		return err
	}

	_, err = statement.Exec(req.Name, req.Age, id)
	if err != nil {
		return err
	}

	return nil
}
