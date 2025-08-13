package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var db *sql.DB

func InitDb() (*sql.DB, error) {
	var err error

	// Initiate secrets and credentials
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	user := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	database := os.Getenv("DATABASE")
	configDatabase := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, database)

	// Connect to local database
	db, err = sql.Open("mysql", configDatabase)
	if err != nil {
		fmt.Println("Connection to localhost database is failed")
		return nil, err
	}

	// Ping the database and error handling
	err = db.Ping()
	if err != nil {
		fmt.Println("Ping to localhost database is failed")
		return nil, err
	}

	// fmt.Println("Connection success")

	return db, nil
}
