package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var DB *sql.DB

func InitDB() {
	var err error

	err = godotenv.Load()

	if err != nil {
		log.Fatal("Error while loading environment variables")
	}

	DB, err = sql.Open("mysql", dbConfig())

	if err != nil {
		log.Fatal("Database Open Error")
	}

	log.Println("Connected To Database")
	createUserTable()
	createBlogTable()
}

func createUserTable() {
	query := "CREATE TABLE IF NOT EXISTS user(id CHAR(36) PRIMARY KEY, email VARCHAR(320), username VARCHAR(100), password VARCHAR(255), joined_on DATE)"

	ctx, cancelFunc := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancelFunc()

	_, err := DB.ExecContext(ctx, query)

	if err != nil {
		log.Fatalf("Error Creating User Table: %s", err)
		return
	}

	log.Println("User Table Has Been Created")
}

func createBlogTable() {
	query := "CREATE TABLE IF NOT EXISTS blog(id CHAR(36) PRIMARY KEY, author_id CHAR(36), title VARCHAR(100), content VARCHAR(255), posted_on DATE, FOREIGN KEY (author_id) REFERENCES user(id))"

	ctx, cancelFunc := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancelFunc()

	_, err := DB.ExecContext(ctx, query)

	if err != nil {
		log.Fatalf("Error Creating Blog Table: %s", err)
		return
	}

	log.Println("Blog Table Has Been Created")
}

func dbConfig() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", os.Getenv("DATABASE_USERNAME"), os.Getenv("DATABASE_PASSWORD"), os.Getenv("DATABASE_HOST"), os.Getenv("DATABASE_NAME"))
}
