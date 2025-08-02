package database

import (
	"database/sql"
	"fmt"
	"github.com/thetestcoder/todo-app/internals/config"
)

func Connect() *sql.DB {
	dsn := config.LoadDbConfig()
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
