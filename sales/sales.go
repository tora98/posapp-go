package sales

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"time"

	"github.com/tora98/posapp-go/utils"
)

type sales struct {
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
		result := utils.GetInput("/sales >")
		if result == "" {
			fmt.Println("Please Enter a Command!")
		}
		command = result

		if slices.Contains(salesCommands, command) {
			switch command {
			case "add":
				var sale sales
				sale.date = time.Now().Format("02, Jan 2006")

				result := utils.GetInput("ID: ")
				if result == "" {
					fmt.Println("Please Enter Product ID :")
					break
				}
				sale.productID = result

				result = utils.GetInput("Quantity: ")
				if result == "" {
					fmt.Println("Please Enter Quantity :")
					break
				}
				sale.productID = result

				result = utils.GetInput("Price: ")
				if result == "" {
					fmt.Println("Please Enter Price :")
					break
				}
				sale.productID = result

				err := sale.valueCheck()
				if err != nil {
					command = "add"
				}

				err = sale.AddSales(db)
				if err != nil {
					return err
				}

			case "delete":
				var sale_id int
				result = utils.GetInput("Id :")
				if result == "" {
					fmt.Println("Invalid Id!")
					break
				}
				sale_id, err := strconv.Atoi(result)
				if err != nil {
					return err
				}

				stmt, err := db.PrepareContext(ctx, "DELETE FROM sales WHERE date = ?, sale_id = ?")
				if err != nil {
					return fmt.Errorf("failed to delete sale: %w", err)
				}
				defer stmt.Close()

				result, err := stmt.ExecContext(ctx, time.Now().Format("02, Jan 2006"), sale_id)
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
				defer stmt.Close()

				rows, err := stmt.QueryContext(ctx, time.Now().Format("02, Jan 2006"))
				if err != nil {
					return fmt.Errorf("failed to execute query: %w", err)
				}
				defer rows.Close()

				for rows.Next() {
					var sale_id int
					var sale sales

					err = rows.Scan(
						&sale_id,
						&sale.date,
						&sale.productID,
						&sale.quantity,
						&sale.price,
					)
					if err != nil {
						return fmt.Errorf("failed to scan rows: %w", err)
					}
					fmt.Printf(
						"%d %s %s %d %d",
						sale_id,
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
				defer stmt.Close()

				rows, err := stmt.QueryContext(ctx)
				if err != nil {
					return fmt.Errorf("failed to execute query: %w", err)
				}
				defer rows.Close()

				for rows.Next() {
					var sale_id int
					var sale sales

					err = rows.Scan(
						&sale_id,
						&sale.date,
						&sale.productID,
						&sale.quantity,
						&sale.price,
					)
					if err != nil {
						return fmt.Errorf("failed to scan rows: %w", err)
					}
					fmt.Printf(
						"%d %s %s %d %d",
						sale_id,
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

// zero value check for sales struct
func (s *sales) valueCheck() error {
	switch {
	case s.productID == "":
		return errors.New("empty product id")
	case s.quantity == 0:
		return errors.New("empty quantity")
	case s.price == 0:
		return errors.New("empty price")
	default:
		return nil
	}
}

// funtion for adding daily sales
func (sale *sales) AddSales(db *sql.DB) error {
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
	defer stmt.Close()

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
