package main

import (
	"fmt"
	"slices"

	"github.com/tora98/posapp-go/products"
	"github.com/tora98/posapp-go/purchases"
	"github.com/tora98/posapp-go/sales"
	"github.com/tora98/posapp-go/utils"

	_ "github.com/mattn/go-sqlite3"
)

// main function
func main() {
	db, err := utils.ConnectDb()
	if err != nil {
		fmt.Println(err)
	}

	mainCommands := []string{"--?", "--help", "sales", "products", "purchases", "exit"}

	var command string
	for command != "exit" {
		result := utils.GetInput("main>")
		command = result

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
