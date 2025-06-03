package utils

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// helper function for connecting to database
func ConnectDb() (*sql.DB, error) {
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
