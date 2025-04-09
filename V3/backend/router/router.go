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
	r.HandleFunc("/stock/in",db.StockIn).Methods("POST")	
	r.HandleFunc("/stock/out",db.StockOut).Methods("POST")	
	r.HandleFunc("/stock/sold",db.StockSold).Methods("POST")	
	r.HandleFunc("/stock/receivedsoldreport",db.ProductReceivedSoldReport).Methods("POST")	
	r.HandleFunc("/stock/productquantity",db.TotalProductQuantityReport).Methods("POST")	

	//Supplier routes
	r.HandleFunc("/supplier/insert",db.InsertSupplier).Methods("POST")
	r.HandleFunc("/supplier/login",db.LoginSupplier).Methods("POST")

	// Inventory routes
	r.HandleFunc("/inventory/add", db.AddProduct).Methods("POST")
	r.HandleFunc("/inventory/remove", db.RemoveProduct).Methods("POST")
	r.HandleFunc("/inventory/data/product", db.GetProductData).Methods("POST")
	r.HandleFunc("/inventory/data/store", db.GetStoreData).Methods("POST")
	r.HandleFunc("/inventory/data/supplier", db.GetSupplierData).Methods("POST")

	return r
}
