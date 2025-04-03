package main

import (
	"V1/db"
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)
func getProductInput() (*db.Product, error) {
        prod := new(db.Product)
        reader := bufio.NewReader(os.Stdin)

        fmt.Print("Enter Name: ")
        name, err := reader.ReadString('\n')
        if err != nil {
                return nil, fmt.Errorf("error reading name: %v", err)
        }
        prod.Name = strings.TrimSpace(name)

        fmt.Print("Enter Description: ")
        desc, err := reader.ReadString('\n')
        if err != nil {
                return nil, fmt.Errorf("error reading description: %v", err)
        }
        prod.Desc = strings.TrimSpace(desc)

        fmt.Print("Enter Price: ")
        priceStr, err := reader.ReadString('\n')
        if err != nil {
                return nil, fmt.Errorf("error reading price: %v", err)
        }
        priceStr = strings.TrimSpace(priceStr)
        prod.Price, err = strconv.ParseFloat(priceStr, 64)
        if err != nil {
                return nil, fmt.Errorf("invalid price: %v", err)
        }

       

        return prod, nil
}
func main() {
    fmt.Println("Hello")
    db.CreateDBConnection()
    defer db.CloseDBConnection()
fmt.Println("Hello")
    db.CreateProductsTable()
    db.CreateStockMovementsTable()
fmt.Println("Hello")
    fmt.Println("System is running!")
    opt := -1
    for opt != 0 {
        fmt.Println("0-End Program")
        fmt.Println("1-Insert Products")
        fmt.Println("2-View Products")
        fmt.Println("3-Update Product")
        fmt.Println("4-Delete Product")
        fmt.Println("5-Stock In")
        fmt.Println("6-Stock Sold")
        fmt.Println("7-Manual Removal")
        fmt.Print("Enter Option:")
        _, err := fmt.Scanln(&opt)
        if err != nil {
            log.Printf("Error reading input: %v", err)
            continue;
        }

        if opt == 1 {
            prod, err := getProductInput()
            if err != nil {
                    fmt.Println("Error in taking input!Try Again")
            }else{
                reader := bufio.NewReader(os.Stdin)
                fmt.Print("Enter Quantity: ")
                quantityStr, err := reader.ReadString('\n')
                if err != nil {
                        fmt.Println("Error Reading quantity!")
                }else{
                    quantityStr = strings.TrimSpace(quantityStr)
                    prod.Quantity, err = strconv.Atoi(quantityStr)
                    if err != nil {
                            fmt.Println("invalid quantity!")
                    }else{
                       db.InsertProduct(*prod)
                    }
                } 
            }
        } else if opt == 2 {
            db.ViewProducts()
        } else if opt==3{
          var id int
          fmt.Print("Enter Product ID:")
          fmt.Scanln(&id)
             prod, err := getProductInput()
        if err != nil {
                fmt.Println("Error in taking input!Try Again")
        }else{
           db.UpdateProduct(*prod,id)
        } 
        }else if opt==4{
             var id int
          fmt.Print("Enter Product ID:")
          fmt.Scanln(&id)
          db.DeleteProduct(id)
        }else if opt==5{
             var id int
          fmt.Print("Enter Product ID:")
          fmt.Scanln(&id)
             reader := bufio.NewReader(os.Stdin)
                fmt.Print("Enter Quantity: ")
                quantityStr, err := reader.ReadString('\n')
                if err != nil {
                        fmt.Println("Error Reading quantity!")
                }else{
                    quantityStr = strings.TrimSpace(quantityStr)
                    Quantity, err := strconv.Atoi(quantityStr)
                    if err != nil {
                            fmt.Println("invalid quantity!")
                    }else{
                       db.StockIn(Quantity,id)
                    }
                } 
        }else if opt==6{
             var id int
          fmt.Print("Enter Product ID:")
          fmt.Scanln(&id)
             reader := bufio.NewReader(os.Stdin)
                fmt.Print("Enter Quantity: ")
                quantityStr, err := reader.ReadString('\n')
                if err != nil {
                        fmt.Println("Error Reading quantity!")
                }else{
                    quantityStr = strings.TrimSpace(quantityStr)
                    Quantity, err := strconv.Atoi(quantityStr)
                    if err != nil {
                            fmt.Println("invalid quantity!")
                    }else{
                       db.StockSold(Quantity,id)
                    }
                } 
        }else if opt==7{
             var id int
          fmt.Print("Enter Product ID:")
          fmt.Scanln(&id)
             reader := bufio.NewReader(os.Stdin)
                fmt.Print("Enter Quantity: ")
                quantityStr, err := reader.ReadString('\n')
                if err != nil {
                        fmt.Println("Error Reading quantity!")
                }else{
                    quantityStr = strings.TrimSpace(quantityStr)
                    Quantity, err := strconv.Atoi(quantityStr)
                    if err != nil {
                            fmt.Println("invalid quantity!")
                    }else{
                       db.StockOut(Quantity,id)
                    }
                } 
        }
    }
}