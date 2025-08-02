package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/thetestcoder/todo-app/internals/server"
)

func main() {
	server.StartServer(server.CreateRoute())
}
