package main

import (
	"context"
	"database/sql"
	"fmt"
	"slices"
	"time"
	//tea "github.com/charmbracelet/bubbletea"
)

type products struct {
	productID         string
	productName       string
	manufacturer      string
	pricePerUnit      int
	pricePerPackaging int
	state             bool
}

type sales struct {
	date      string
	productID string
	quantity  int
	price     int
}

type purchases struct {
	date         string
	productID    string
	productName  string
	manufacturer string
	quantity     int
	price        int
	supplier     string
}

// main function
func main() {
	//	db, err := connectDb(); if err != nil {
	//		fmt.Println(err)
	//	}

	commands := []string{"--?", "--help", "sales", "products", "purchases", "exit"}

	var command string
	for command != "exit" {
		fmt.Print("> ")
		_, err := fmt.Scanln(&command)
		if err != nil {
			fmt.Println(err)
		}
		if slices.Contains(commands, command) {
			switch command {
			case "":
				help()
			case "--?":
				help()
			case "sales":
				fmt.Println("Sales")
			case "products":
				fmt.Println("Products")
			case "purchases":
				fmt.Println("Purchases")
			}
		} else {
			fmt.Println("Unknown Command!")
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

//helper function for checking if product is in database
//func checkExists(item interface{}) bool {
//	return false
//}

// function that prints the help/options
func help() {
	fmt.Println("POSapp Options!")
	fmt.Println("-------------------------------------")
	fmt.Println("sales      --Enter sales mode for interacting with daily sales!")
	fmt.Println("    add      --Adds a daily sale. Accepts 3 arguments separated by spaces(product_name, quantity, price) ")
	fmt.Println("    delete      --Deletes a daily sale. Accepts 1 argument(product_id) ")
	fmt.Println("    list      --Lists daily sales. Accepts 0 arguments")
	fmt.Println("    listAll      --Lists all daily sales from database. Accepts 0 arguments")
	fmt.Println("purchases      --Enter purchases mode for interacting with purchases!")
	fmt.Println(`    add      --Adds a purchase on the current date. Accepts 6 arguments separated by spaces
																(product_id, product_name, manufacturer, price_per_unit, prce_per_packaging, state)`)
	fmt.Println("    list      --Lists all active products!. Accepts 0 arguments")
	fmt.Println("products      --Enter products mode for interacting with products!")
	fmt.Println(`    add      --Adds a new product. Accepts 6 arguments separated by spaces)
																(product_id, product_name, manufacturer, quantity, price, supplier)`)
	fmt.Println(`    list      --Lists active products from database. Accepts 0 arguments`)
	fmt.Println(`    delete      --Disables an active products from database. Accepts 1 argument (product_id)`)
	fmt.Println(`    enable      --Disables an active products from database. Accepts 1 argument (product_id)`)
}

// function for adding purchase exoenses
func (purchase *purchases) addPurchase(db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stmt, err := db.PrepareContext(ctx, `INSERT INTO purchases(
		date,
		product_id,
		product_nameame,
		manufacturer,
		quantity,
		price,
		supplier
		) VALUES(?, ?, ?, ?, ?, ?, ?)`,
	)
	if err != nil {
		return fmt.Errorf("faile to prepare query statement: %w", err)
	}

	defer stmt.Close()

	result, err := stmt.ExecContext(
		ctx,
		purchase.date,
		purchase.productID,
		purchase.productName,
		purchase.manufacturer,
		purchase.quantity,
		purchase.price,
		purchase.supplier,
	)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to query rows affected: %w", err)
	}
	fmt.Printf("Purchase Added Succeddfully! rowsaffected=?\n", affectedRows)
	return nil
}

// funtion for adding daily sales
func (sale *sales) addSales(db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stmt, err := db.PrepareContext(ctx, `INSERT INTO sales(
		date,
		product_id,
		quantity,
		price) VALUES(?, ?, ?, ?)`,
	)
	if err != nil {
		return fmt.Errorf("failed to prepare query statement: %w", err)
	}

	result, err := stmt.ExecContext(
		ctx,
		sale.date,
		sale.productID,
		sale.quantity,
		sale.price,
	)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to query rows affected: %w", err)
	}

	fmt.Printf("Sales Added Succesfully! rowsaffected=?\n", rowsAffected)
	return nil
}

// function for adding a product into products table
func (product *products) addProduct(db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stmt, err := db.PrepareContext(ctx, `INSERT INTO products(
		product_id,
		product_name,
		manufacturer,
		price_per_unit,
		prce_per_packaging,
		status) VALUES(?, ?, ?, ?, ?, ?)`,
	)
	if err != nil {
		return fmt.Errorf("failed to prepare query: %w", err)
	}

	defer stmt.Close()

	result, err := stmt.ExecContext(
		ctx,
		product.productID,
		product.productName,
		product.manufacturer,
		product.pricePerUnit,
		product.pricePerPackaging,
		product.state,
	)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to query insert id: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to query rows affected %w", err)
	}

	fmt.Printf("Product Saved Successfully! id=?, rowsaffected=?\n", id, rowsAffected)
	return nil
}
