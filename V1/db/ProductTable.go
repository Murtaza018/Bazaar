package db

import (
	"fmt"
	"log"
)
type Product struct {
    Name string
    Price  float64
	Desc string
	Quantity int
}
// CreateProductsTable creates the Products table if it doesn't exist
func CreateProductsTable() {
	_, err := Db.Exec(`
		CREATE TABLE IF NOT EXISTS Products (
			product_id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT,
			description TEXT,
			price REAL,
			quantity INTEGER
		);
	`)
	if err != nil {
		log.Fatal("Error creating Products table:", err)
	}
	fmt.Println("Products table created (or already exists).")
}
func ViewProducts(){
	rows, err := Db.Query("SELECT * FROM Products") 
        if err != nil {
                log.Fatal("Error querying Products table:", err)
        }
        defer rows.Close()

        // Get column names (optional)
        columns, err := rows.Columns()
        if err != nil {
                log.Fatal("Error getting column names:", err)
        }

        // Iterate over rows
        for rows.Next() {
                // Create a slice of interface{} to hold the row values
                values := make([]interface{}, len(columns))
                valuePtrs := make([]interface{}, len(columns))
                for i := range values {
                        valuePtrs[i] = &values[i]
                }

                // Scan the row into the values slice
                err := rows.Scan(valuePtrs...)
                if err != nil {
                        log.Fatal("Error scanning row:", err)
                }

                // Display the row values
                for i, col := range columns {
                        fmt.Printf("%s: %v, ", col, values[i])
                }
                fmt.Println() // Newline for each row
        }

        // Check for errors during row iteration
        if err := rows.Err(); err != nil {
                log.Fatal("Error during row iteration:", err)
        }

}
func InsertProduct(prod Product) {
    _, err := Db.Exec(`INSERT INTO Products (name, description, price, quantity) VALUES (?, ?, ?, ?)`, prod.Name, prod.Desc, prod.Price, prod.Quantity)
    if err != nil {
        log.Fatal("Error inserting product:", err)
    }
    fmt.Println("Product Inserted!")
}
func UpdateProduct(prod Product,id int) {
    _, err := Db.Exec(`Update Products set name=?,description=?,price=? where product_id=?`, prod.Name, prod.Desc, prod.Price,id)
    if err != nil {
        log.Fatal("Error updating product:", err)
    }
    fmt.Println("Product Updated!")
}
func DeleteProduct(id int) {
    _, err := Db.Exec(`Delete from Products where product_id=?`,id)
    if err != nil {
        log.Fatal("Error deleting product:", err)
    }
    fmt.Println("Product Deleted!")
}