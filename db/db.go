package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func Connect() {
	var err error
	DB, err = sql.Open("sqlite3", "./todo.db")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to sqlite database...")
}

func CreateTable() {
	table := "CREATE TABLE IF NOT EXISTS todo (id INTEGER PRIMARY KEY AUTOINCREMENT, item TEXT, isCompleted INTEGER)"
	_, err := DB.Exec(table)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Table created...")
}
