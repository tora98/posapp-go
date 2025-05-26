package main

import (
	"database/sql"
	"fmt"
	"posapp-go/products"
	"slices"

	//"posapp-go/purchases"
	//"posapp-go/sales"

	_ "github.com/mattn/go-sqlite3"
	//tea "github.com/charmbracelet/bubbletea"
)

// main function
func main() {
	db, err := connectDb()
	if err != nil {
		fmt.Println(err)
	}

	mainCommands := []string{"--?", "--help", "sales", "products", "purchases", "exit"}

	var command string
	for command != "exit" {
		fmt.Print("main > ")
		_, err := fmt.Scanln(&command)
		if err != nil {
			fmt.Println(err)
		}
		if slices.Contains(mainCommands, command) {
			switch command {
			case "--?":
				help()
			case "--help":
				help()
			case "sales":
				fmt.Println("Sales")
			case "products":
				err := products.Menu(db)
				if err != nil {
					fmt.Printf("Error: %v", err)
				}
			case "purchases":
				fmt.Println("Purchases")
			}
		}
	}
}

// helper function for connecting to database
func connectDb() (*sql.DB, error) {
	const database = "./posapp.db"
	db, err := sql.Open("sqlite3", database)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

// function that prints the help/options
func help() {
	fmt.Println("POSapp Options!")
	fmt.Println("-----------------------------------------------------------------------------------------------------------------------------------")
	fmt.Println("sales      	--Enter sales mode for interacting with daily sales!")
	fmt.Println("    add      	 	--Adds a daily sale. Accepts 3 additional arguments separated by spaces(product_name, quantity, price) ")
	fmt.Println("    delete      	--Deletes a daily sale. Accepts 1 additional argument(product_id) ")
	fmt.Println("    list      	 	--Lists daily sales. Accepts 0 arguments")
	fmt.Println("    listAll     	--Lists all daily sales from database. Accepts 0 arguments")
	fmt.Println("purchases      --Enter purchases mode for interacting with purchases!")
	fmt.Println(`    add      		--Adds a purchase on the current date. Accepts 6 additional arguments separated by spaces
																(product_id, product_name, manufacturer, price_per_unit, prce_per_packaging, state)`)
	fmt.Println("    list      		--Lists all active products!. Accepts 0 arguments")
	fmt.Println("products       --Enter products mode for interacting with products!")
	fmt.Println(`    add      		--Adds a new product. Accepts 6 additional arguments separated by spaces)
																(product_id, product_name, manufacturer, quantity, price, supplier)`)
	fmt.Println(`    list      		--Lists active products from database. Accepts 0 arguments`)
	fmt.Println(`    delete      	--Disables an active products from database. Accepts 1 additional argument (product_id)`)
	fmt.Println(`    enable      	--Disables an active products from database. Accepts 1 additional argument (product_id)`)
}
