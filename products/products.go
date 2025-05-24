package products

import (
	"context"
	"database/sql"
	"fmt"
	"slices"
	"time"
)

type products struct {
	productID         string
	productName       string
	manufacturer      string
	pricePerUnit      int
	pricePerPackaging int
	state             bool
}

// Products menu
func ProductMenu(db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	commands := []string{
		"add",
		"list",
		"delete",
		"enable",
		"quit",
	}

	var command string
	for command != "quit" {
		fmt.Print("/products >")

		_, err := fmt.Scanln(&command)
		if err != nil {
			fmt.Println(err)
		}

		if slices.Contains(commands, command) {
			switch command {
			case "add":
				var product products

				fmt.Print("ID: ")
				_, err := fmt.Scanln(&product.productID)
				if err != nil {
					return err
				}
				fmt.Print("Product Name: ")
				_, err = fmt.Scanln(&product.productName)
				if err != nil {
					return err
				}
				fmt.Print("Manufacturer: ")
				_, err = fmt.Scanln(&product.manufacturer)
				if err != nil {
					return err
				}
				fmt.Print("Price Per Unit: ")
				_, err = fmt.Scanln(&product.pricePerUnit)
				if err != nil {
					return err
				}
				fmt.Print("Price Per Packaging: ")
				_, err = fmt.Scanln(&product.pricePerPackaging)
				if err != nil {
					return err
				}
				fmt.Print("State: ")
				_, err = fmt.Scanln(&product.state)
				if err != nil {
					return err
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
					var id string
					var productName string
					var manufacturer string
					var pricePerUnit int
					var pricePerPackaging int
					var state bool

					err = rows.Scan(&id, &productName, &manufacturer, &pricePerUnit, &pricePerPackaging, &state)
					if err != nil {
						return fmt.Errorf("failed to scan rows: %w", err)
					}
					fmt.Printf("%s %s %s %d %d %t", id, productName, manufacturer, pricePerUnit, pricePerPackaging, state)
				}

			case "delete":
				var id string
				fmt.Print("id: ")
				_, err := fmt.Scanln(&id)
				if err != nil {
					return err
				}

				stmt, err := db.PrepareContext(ctx, "UPDATE products SET state = ? WHERE product_id = ?")
				if err != nil {
					return fmt.Errorf("failed to disable product: %w", err)
				}

				result, err := stmt.ExecContext(ctx, false, id)
				if err != nil {
					return fmt.Errorf("failed to execute query: %w", err)
				}
				affectedRows, err := result.RowsAffected()
				if err != nil {
					return fmt.Errorf("failed to query rows affected: %w", err)
				}

				fmt.Printf("Disabled Product Successfully! rowsaffected=%d\n", affectedRows)

			case "enable":
				var id string
				fmt.Print("id: ")
				_, err := fmt.Scanln(&id)
				if err != nil {
					return err
				}

				stmt, err := db.PrepareContext(ctx, "UPDATE products SET state = ? WHERE product_id = ?")
				if err != nil {
					return fmt.Errorf("failed to disable product: %w", err)
				}

				result, err := stmt.ExecContext(ctx, true, id)
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
			fmt.Println("Invalid Command!")
		}
	}
	return nil
}

// function for adding a product into products table
func (product *products) AddProduct(db *sql.DB) error {
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
