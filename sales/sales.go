package sales

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type sales struct {
	date      string
	productID string
	quantity  int
	price     int
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
