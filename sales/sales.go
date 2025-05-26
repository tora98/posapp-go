package sales

import (
	"context"
	"database/sql"
	"fmt"
	"slices"
	"time"
)

type sales struct {
	id        int
	date      string
	productID string
	quantity  int
	price     int
}

// Sales menu
func Menu(db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	salesCommands := []string{
		"add",
		"delete",
		"list",
		"listAll",
		"quit",
	}

	var command string
	for command != "quit" {
		fmt.Print("/sales >")

		_, err := fmt.Scanln(&command)
		if err != nil {
			fmt.Println(err)
		}

		if slices.Contains(salesCommands, command) {
			switch command {
			case "add":
				var sales sales
				sales.date = time.Now().Format("02, Jan 2006")

				fmt.Print("ID: ")
				_, err := fmt.Scanln(&sales.productID)
				if err != nil {
					return err
				}
				fmt.Print("Quantity: ")
				_, err = fmt.Scanln(&sales.quantity)
				if err != nil {
					return err
				}
				fmt.Print("Price: ")
				_, err = fmt.Scanln(&sales.price)
				if err != nil {
					return err
				}
				err = sales.addSales(db)
				if err != nil {
					return err
				}

			case "delete":
				var product_id string
				fmt.Print("id: ")
				_, err := fmt.Scanln(&product_id)
				if err != nil {
					return err
				}

				stmt, err := db.PrepareContext(ctx, "DELETE FROM sales WHERE date = ?, product_id = ?")
				if err != nil {
					return fmt.Errorf("failed to delete sale: %w", err)
				}

				result, err := stmt.ExecContext(ctx, time.Now().Format("02, Jan 2006"), product_id)
				if err != nil {
					return fmt.Errorf("failed to execute query: %w", err)
				}
				affectedRows, err := result.RowsAffected()
				if err != nil {
					return fmt.Errorf("failed to query rows affected: %w", err)
				}

				fmt.Printf("Disabled Product Successfully! rowsaffected=%d\n", affectedRows)

			case "list":
				fmt.Println("=============================================================")

				stmt, err := db.PrepareContext(ctx, "SELECT * FROM sales WHERE date = ?")
				if err != nil {
					return fmt.Errorf("failed to prepare query: %w", err)
				}

				rows, err := stmt.QueryContext(ctx, time.Now().Format("02, Jan 2006"))
				if err != nil {
					return fmt.Errorf("failed to execute query: %w", err)
				}
				defer rows.Close()

				for rows.Next() {
					var sale sales

					err = rows.Scan(
						&sale.date,
						&sale.productID,
						&sale.quantity,
						&sale.price,
					)
					if err != nil {
						return fmt.Errorf("failed to scan rows: %w", err)
					}
					fmt.Printf(
						"%s %s %d %d",
						sale.date,
						sale.productID,
						sale.quantity,
						sale.price,
					)
				}
			case "listAll":
				fmt.Println("=============================================================")

				stmt, err := db.PrepareContext(ctx, "SELECT * FROM sales")
				if err != nil {
					return fmt.Errorf("failed to prepare query: %w", err)
				}

				rows, err := stmt.QueryContext(ctx)
				if err != nil {
					return fmt.Errorf("failed to execute query: %w", err)
				}
				defer rows.Close()

				for rows.Next() {
					var sale sales

					err = rows.Scan(
						&sale.date,
						&sale.productID,
						&sale.quantity,
						&sale.price,
					)
					if err != nil {
						return fmt.Errorf("failed to scan rows: %w", err)
					}
					fmt.Printf(
						"%s %s %d %d",
						sale.date,
						sale.productID,
						sale.quantity,
						sale.price,
					)
				}
			}
		}
	}

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

	fmt.Printf("Sales Added Succesfully! rowsaffected=%d\n", rowsAffected)
	return nil
}
