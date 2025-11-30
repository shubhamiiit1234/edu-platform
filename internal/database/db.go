package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitializeDB(connection string) error {
	db, err := sql.Open("postgres", connection)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}
	DB = db

	err = DB.Ping()
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	fmt.Println("Connected to Database successfully!!!")
	return nil
}

func GetDBInstance() (*sql.DB, error) {
	if DB == nil {
		return nil, fmt.Errorf("database not initialized")
	}
	return DB, nil
}
