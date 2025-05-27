package main

import (
	"context"
	"database/sql"
	"fmt"
	"posapp-go/products"
	"posapp-go/purchases"
	"posapp-go/sales"
	"slices"
	"time"

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
			case "--?", "--help":
				help()
			case "sales":
				err := sales.Menu(db)
				if err != nil {
					fmt.Printf("Error: %v", err)
				}
			case "products":
				err := products.Menu(db)
				if err != nil {
					fmt.Printf("Error: %v", err)
				}
			case "purchases":
				err := purchases.Menu(db)
				if err != nil {
					fmt.Printf("Error: %v", err)
				}
			case "exit":
				fmt.Println("Exiting...")
			default:
				fmt.Println("Invalid command")
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
	err = createTables(db)
	if err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}
	return db, nil
}

// create tables for sales, products and purchases
func createTables(db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS products (
		product_id TEXT PRIMARY KEY,
		product_name TEXT NOT NULL,
		manufacturer TEXT NOT NULL,
		quantity INTEGER NOT NULL,
		price_per_unit REAL NOT NULL,
		price_per_packaging REAL NOT NULL,
		state INTEGER NOT NULL
	)`,
	)
	if err != nil {
		return fmt.Errorf("failed to create products table: %w", err)
	}

	_, err = db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS sales (
		sale_id INTEGER PRIMARY KEY AUTOINCREMENT,
		date TEXT NOT NULL,
		product_id TEXT UNIQUE NOT NULL,
		quantity INTEGER NOT NULL,
		price REAL NOT NULL,
		FOREIGN KEY (product_id) REFERENCES products (product_id))`,
	)
	if err != nil {
		return fmt.Errorf("failed to create sales table: %w", err)
	}

	_, err = db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS purchases (
		purchase_id INTEGER PRIMARY KEY AUTOINCREMENT,
		date TEXT NOT NULL,
		product_id TEXT NOT NULL,
		product_name TEXT NOT NULL,
		manufacturer TEXT NOT NULL,
		quantity INTEGER NOT NULL,
		price REAL NOT NULL,
		supplier TEXT NOT NULL)`,
	)
	if err != nil {
		return fmt.Errorf("failed to create purchases table: %w", err)
	}

	return nil
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
	fmt.Println("products      --Enter purchases mode for interacting with purchases!")
	fmt.Println(`    add      		--Adds a purchase on the current date. Accepts 6 additional arguments separated by spaces
																(product_id, product_name, manufacturer, price_per_unit, prce_per_packaging, state)`)
	fmt.Println("    list      		--Lists all active products!. Accepts 0 arguments")
	fmt.Println("	 delete			--Disables an active product from database. Accepts 1 additional argument (product_id)")
	fmt.Println(" 	 enable 		--Enables an existing but disabled product. Accepts 1 additional argument (product_id)")

	fmt.Println("purchases       --Enter products mode for interacting with products!")
	fmt.Println(`    add      		--Adds a new product. Accepts 6 additional arguments separated by spaces)
																(product_id, product_name, manufacturer, quantity, price, supplier)`)
	fmt.Println("    delete      	--Disables an active products from database. Accepts 1 additional argument (product_id)")
	fmt.Println("    list      		--Lists active products from database. Accepts 0 arguments")
	fmt.Println("    listAll      	--Lists all products from database. Accepts 0 arguments")
}
