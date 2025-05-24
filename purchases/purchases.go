package purchases

import (
	"context"
	"database/sql"
	"fmt"
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
	fmt.Printf("Purchase Added Succeddfully! rowsaffected=%d\n", affectedRows)
	return nil
}
