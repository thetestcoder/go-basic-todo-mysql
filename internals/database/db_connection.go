package database

import (
	"database/sql"
	"fmt"
)

func Connect() *sql.DB {
	dsn := "root:secret@tcp(localhost:3306)/todo_app"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	fmt.Println(
		"Connected to database",
		db.Stats().OpenConnections,
		"connections",
	)
	return db
}

func Close(db *sql.DB) {
	db.Close()
}
