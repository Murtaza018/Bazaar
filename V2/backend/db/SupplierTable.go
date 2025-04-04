package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)
func CreateSupplierTable(Db *sql.DB) {
	_, err := Db.Exec(`
		CREATE TABLE IF NOT EXISTS Supplier (
			supplier_id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			password TEXT NOT NULL,
			location TEXT,
			);
	`)
	if err != nil {
		log.Fatal("Error creating Store table:", err)
	}
	fmt.Println("Store table created (or already exists).")
}
func EncryptPassword(password string) string {
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator

	// Generate a random number between 1 and 9
	randomNumber := rand.Intn(9) + 1
	var encryptedPassword strings.Builder

	// Encrypt each character by adding the random number to its ASCII value
	for _, char := range password {
		encryptedPassword.WriteRune(char + rune(randomNumber))
	}

	// Append the random number at the end
	encryptedPassword.WriteString(strconv.Itoa(randomNumber))

	return encryptedPassword.String()
}

// DecryptPassword decrypts the password by subtracting the appended random number from each character's ASCII value.
func DecryptPassword(encryptedPassword string) string {
	// Extract the last character (which is the random number)
	length := len(encryptedPassword)
	randomNumber, err := strconv.Atoi(string(encryptedPassword[length-1])) // Convert last character to integer
	if err != nil {
		log.Fatal("Error decrypting password:", err)
	}

	// Decrypt the password by subtracting the random number from each character's ASCII value
	var decryptedPassword strings.Builder
	for _, char := range encryptedPassword[:length-1] {
		decryptedPassword.WriteRune(char - rune(randomNumber))
	}

	return decryptedPassword.String()
}

// InsertSupplier inserts a new supplier into the Supplier table
func InsertSupplier(w http.ResponseWriter, r *http.Request) {
	// Parsing input from the request body
	var supplier struct {
		Name     string `json:"name"`
		Password string `json:"password"`
		Location string `json:"location"`
	}

	// Parse the JSON request body
	err := json.NewDecoder(r.Body).Decode(&supplier)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Encrypt the password before inserting into the database
	encryptedPassword := EncryptPassword(supplier.Password)

	// Insert into the database
	_, err = Db.Exec(`INSERT INTO Supplier (name, password, location) VALUES ($1, $2, $3)`,
		supplier.Name, encryptedPassword, supplier.Location)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error inserting supplier: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Supplier inserted!")
}
func LoginSupplier(w http.ResponseWriter, r *http.Request) {
	// Parse the JSON request body to get supplier_id and password
	var requestBody struct {
		SupplierID int    `json:"supplier_id"`
		Password   string `json:"password"`
	}

	// Decode the request body into the requestBody struct
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Query the database to get the encrypted password for the given supplier_id
	var encryptedPassword string
	query := `SELECT password FROM Supplier WHERE supplier_id = $1`
	err = Db.QueryRow(query, requestBody.SupplierID).Scan(&encryptedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Supplier not found", http.StatusNotFound)
		} else {
			log.Println("Error querying supplier:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// Decrypt the password retrieved from the database
	decryptedPassword := DecryptPassword(encryptedPassword)
	
	// Compare the decrypted password with the entered password
	if decryptedPassword != requestBody.Password {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	// Respond with a success message if login is successful
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Login successful!"))
}