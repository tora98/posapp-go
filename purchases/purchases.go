package purchases

import (
	"context"
	"database/sql"
	"fmt"
	"slices"
	"time"
)

type purchases struct {
	date         string
	productID    string
	productName  string
	manufacturer string
	quantity     int
	price        int
	supplier     string
}

// Purchases menu
func Menu(db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	purchasesCommands := []string{
		"add",
		"delete",
		"list",
		"listAll",
		"quit",
	}

	var command string
	for command != "quit" {
		fmt.Print("/purchases >")

		_, err := fmt.Scanln(&command)
		if err != nil {
			fmt.Println(err)
		}

		if slices.Contains(purchasesCommands, command) {
			switch command {
			case "add":
				var purchase purchases
				purchase.date = time.Now().Format("02, Jan 2006")

				fmt.Print("ID: ")
				_, err := fmt.Scanln(&purchase.productID)
				if err != nil {
					return err
				}
				fmt.Print("Quantity: ")
				_, err = fmt.Scanln(&purchase.productName)
				if err != nil {
					return err
				}
				fmt.Print("Price: ")
				_, err = fmt.Scanln(&purchase.manufacturer)
				if err != nil {
					return err
				}
				fmt.Print("Price: ")
				_, err = fmt.Scanln(&purchase.quantity)
				if err != nil {
					return err
				}
				fmt.Print("Price: ")
				_, err = fmt.Scanln(&purchase.price)
				if err != nil {
					return err
				}
				fmt.Print("Price: ")
				_, err = fmt.Scanln(&purchase.supplier)
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

				stmt, err := db.PrepareContext(ctx, "DELETE FROM purchases WHERE date = ?, product_id = ?")
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

// function for adding purchase exoenses
func (purchase *purchases) AddPurchase(db *sql.DB) error {
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
	fmt.Printf("Purchase Added Succeddfully! rowsaffected=%d\n", affectedRows)
	return nil
}
