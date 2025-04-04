package router

import (
	"backend/db" // Replace with your module path
	"net/http"

	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	r := mux.NewRouter()

	// Public endpoints
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("API is running!"))
	}).Methods("GET")

	// Store routes
	r.HandleFunc("/store/register", db.InsertStore).Methods("POST")
	r.HandleFunc("/store/login", db.LoginStore).Methods("POST")

	// Product routes
	r.HandleFunc("/products/view", db.ViewProducts).Methods("GET")
	r.HandleFunc("/products/insert", db.InsertProduct).Methods("POST")
	r.HandleFunc("/products/update", db.UpdateProduct).Methods("POST")
	r.HandleFunc("/products/delete", db.DeleteProduct).Methods("POST")
	
	//Stock Movement routes
	r.HandleFunc("/stock/stockin",db.StockIn).Methods("POST")	
	r.HandleFunc("/stock/stockout",db.StockOut).Methods("POST")	
	r.HandleFunc("/stock/stocksold",db.StockSold).Methods("POST")	

	//Supplier routes
	r.HandleFunc("/supplier/insert",db.InsertSupplier).Methods("POST")
	r.HandleFunc("/supplier/login",db.LoginSupplier).Methods("POST")
	
	// Inventory routes
	r.HandleFunc("/inventory", db.UpdateInventory).Methods("POST")

	return r
}
