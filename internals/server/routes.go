package server

import (
	"github.com/gorilla/mux"
	"github.com/thetestcoder/todo-app/internals/database"
	"github.com/thetestcoder/todo-app/internals/handler"
)

// CreateRoute it will create routes for todo application
func CreateRoute() *mux.Router {
	router := mux.NewRouter()
	initializeTodoRoutes(router)
	return router
}

func initializeTodoRoutes(router *mux.Router) {
	todoHandler := handler.NewTodoHandler(database.Connect())
	router.HandleFunc("/create", todoHandler.CreateTodo).Methods("POST")
	router.HandleFunc("/list", todoHandler.GetTodos).Methods("GET")
	router.HandleFunc("/update/{id}", todoHandler.UpdateTodo).Methods("PUT")
	router.HandleFunc("/delete/{id}", todoHandler.DeleteTodo).Methods("DELETE")
}
