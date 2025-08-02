package server

import (
	"github.com/gorilla/mux"
	"github.com/thetestcoder/todo-app/internals/handler"
)

// CreateRoute it will create routes for todo application
func CreateRoute() *mux.Router {
	router := mux.NewRouter()
	initializeTodoRoutes(router)
	return router
}

func initializeTodoRoutes(router *mux.Router) {
	router.HandleFunc("/create", handler.CreateTodo).Methods("POST")
	router.HandleFunc("/list", handler.GetTodos).Methods("GET")
	router.HandleFunc("/update/{id}", handler.UpdateTodo).Methods("PUT")
	router.HandleFunc("/delete/{id}", handler.DeleteTodo).Methods("DELETE")
}
