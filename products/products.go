package products

import (
	"context"
	"database/sql"
	"fmt"
	"slices"
	"strconv"
	"time"

	"github.com/tora98/posapp-go/utils"
)

type product struct {
	productID         string
	productName       string
	manufacturer      string
	pricePerUnit      int
	pricePerPackaging int
	state             bool
}

// Products menu
func Menu(db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	productCommands := []string{
		"add",
		"list",
		"delete",
		"enable",
		"quit",
	}

	var command string

	for command != "quit" {
		result := utils.GetInput("/products >")
		if result == "" {
			fmt.Println("Please Enter a Command!")
		}
		command = result

		if slices.Contains(productCommands, command) {
			switch command {
			case "add":
				var product product
				product.state = true

				result := utils.GetInput("ID :")
				if result == "" {
					fmt.Println("Please Enter ID!")
					break
				}
				product.productID = result

				result = utils.GetInput("Product Name :")
				if result == "" {
					fmt.Println("Please Enter Product Name!")
					break
				}
				product.productName = result

				result = utils.GetInput("Manufacturer :")
				if result == "" {
					fmt.Println("Please Enter Manufacturer!")
					break
				}
				product.manufacturer = result

				result = utils.GetInput("Price Per Unit :")
				if result == "" {
					fmt.Println("Please Enter a Price per Unit!")
					break
				}
				var err error
				product.pricePerUnit, err = strconv.Atoi(result)
				if err != nil {
					fmt.Println("Please Enter a Valid Price!")
					break
				}

				result = utils.GetInput("Price Per Packaging :")
				if result == "" {
					fmt.Println("Please Enter a Price per Packaging!")
					break
				}
				product.pricePerPackaging, err = strconv.Atoi(result)
				if err != nil {
					fmt.Println("Please Enter a Valid Price!")
					break
				}

				check, err := product.checkExists(db)
				if err != nil {
					return err
				}
				if check {
					fmt.Println("Product already exists")
					command = "add" // go back to add command :-|
				}

				err = product.AddProduct(db)
				if err != nil {
					return err
				}

			case "list":
				fmt.Println("=============================================================")

				stmt, err := db.PrepareContext(ctx, "SELECT * FROM products WHERE state = ?")
				if err != nil {
					return fmt.Errorf("failed to prepare query: %w", err)
				}

				rows, err := stmt.QueryContext(ctx, true)
				if err != nil {
					return fmt.Errorf("failed to execute query: %w", err)
				}
				defer rows.Close()

				for rows.Next() {
					var product product

					err = rows.Scan(
						&product.productID,
						&product.productName,
						&product.manufacturer,
						&product.pricePerUnit,
						&product.pricePerPackaging,
						&product.state,
					)
					if err != nil {
						return fmt.Errorf("failed to scan rows: %w", err)
					}
					fmt.Printf(
						"%s %s %s %d %d %t",
						product.productID,
						product.productName,
						product.manufacturer, product.pricePerUnit,
						product.pricePerPackaging,
						product.state,
					)
				}

			case "delete":
				var product_id string
				result = utils.GetInput("ID: ")
				if result == "" {
					fmt.Println("Invalid ID!")
					break
				}
				product_id = result

				stmt, err := db.PrepareContext(ctx, "UPDATE products SET state = ? WHERE product_id = ?")
				if err != nil {
					return fmt.Errorf("failed to disable product: %w", err)
				}

				result, err := stmt.ExecContext(ctx, false, product_id)
				if err != nil {
					return fmt.Errorf("failed to execute query: %w", err)
				}
				affectedRows, err := result.RowsAffected()
				if err != nil {
					return fmt.Errorf("failed to query rows affected: %w", err)
				}

				fmt.Printf("Disabled Product Successfully! rowsaffected=%d\n", affectedRows)

			case "enable":
				var product_id string
				result = utils.GetInput("ID: ")
				if result == "" {
					fmt.Println("Invalid ID!")
					break
				}
				product_id = result

				stmt, err := db.PrepareContext(ctx, "UPDATE products SET state = ? WHERE product_id = ?")
				if err != nil {
					return fmt.Errorf("failed to disable product: %w", err)
				}

				result, err := stmt.ExecContext(ctx, true, product_id)
				if err != nil {
					return fmt.Errorf("failed to execute query: %w", err)
				}
				affectedRows, err := result.RowsAffected()
				if err != nil {
					return fmt.Errorf("failed to query rows affected: %w", err)
				}

				fmt.Printf("Enabled Product Successfully! rowsaffected=%d\n", affectedRows)
			}
		} else {
			fmt.Println("Not a Valid Command!")
		}
	}
	return nil
}

// function for adding a product into products table
func (product *product) AddProduct(db *sql.DB) error {
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

	fmt.Printf("Product Saved Successfully! id=%d, rowsaffected=%d\n", id, rowsAffected)
	return nil
}

func (product *product) checkExists(db *sql.DB) (bool, error) {
	ctx, calcel := context.WithTimeout(context.Background(), 5*time.Second)
	defer calcel()

	stmt, err := db.PrepareContext(ctx, "SELECT product_id FROM products WHERE product_id = ?")
	if err != nil {
		return false, fmt.Errorf("failed to prepare query: %w", err)
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, product.productID)

	var productID string
	err = row.Scan(&productID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("failed to scan row: %w", err)
	}

	return true, nil
}
