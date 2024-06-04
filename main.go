package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"gotest/database"
	"gotest/model"
	"gotest/service"
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

		customers, err := service.GetAllCustomer(db)
		if err != nil {
			message := []byte(`{"message": "query error"}`)
			setJsonResp(message, http.StatusInternalServerError, res)
			return
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

		err = service.CreateCustomer(customer, db)
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
		customer, err := service.GetCustomerById(id, db)
		if err != nil {
			message := []byte(`{"message": "Error GetCustomerById"}`)
			setJsonResp(message, http.StatusInternalServerError, res)
			return
		}

		if customer.Name == "" {
			message := []byte(`{"message": "not found"}`)
			setJsonResp(message, http.StatusNotFound, res)
			return
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
		err := service.DeleteCustomer(id, db)
		if err != nil {
			message := []byte(`{"message": "Error delete data"}`)
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

		customer, err := service.GetCustomerById(id, db)
		if err != nil {
			message := []byte(`{"message": "Error GetCustomerById"}`)
			setJsonResp(message, http.StatusInternalServerError, res)
			return
		}

		if customer.Name == "" {
			message := []byte(`{"message": "not found"}`)
			setJsonResp(message, http.StatusNotFound, res)
			return
		}

		fmt.Printf("name : %v | age : %v\n", updateCustomer.Name, updateCustomer.Age)

		err = service.UpdateCustomer(id, updateCustomer, db)
		if err != nil {
			message := []byte(`{"message": "Error UpdateCustomer"}`)
			setJsonResp(message, http.StatusInternalServerError, res)
			return
		}

		message := []byte(`{"message": "SUCCESS"}`)
		setJsonResp(message, http.StatusOK, res)
		return
	}

}

func main() {
	database.InitDatabase() //new database every time
	fmt.Println("Server start port: 8080")
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
