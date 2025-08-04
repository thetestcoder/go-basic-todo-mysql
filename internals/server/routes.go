package server

import (
	"github.com/gorilla/mux"
	"github.com/thetestcoder/todo-app/internals/database"
	"github.com/thetestcoder/todo-app/internals/handler"
	"github.com/thetestcoder/todo-app/internals/middleware"
	"github.com/thetestcoder/todo-app/internals/repository"
	"github.com/thetestcoder/todo-app/internals/service"
	"github.com/thetestcoder/todo-app/internals/validator"
)

// CreateRoute it will create routes for todo application
func CreateRoute() *mux.Router {
	router := mux.NewRouter()
	initializeTodoRoutes(router)
	return router
}

func initializeTodoRoutes(router *mux.Router) {

	todoRepository := repository.NewSQLTodoRepository(database.Connect())
	todoService := service.NewTodoService(todoRepository, validator.NewTodoValidator())
	todoHandler := handler.NewTodoHandler(todoService)

	router.Use(middleware.RequestLogger)

	router.HandleFunc("/create", todoHandler.CreateTodo).Methods("POST")
	router.HandleFunc("/list", todoHandler.GetTodos).Methods("GET")
	router.HandleFunc("/update/{id}", todoHandler.UpdateTodo).Methods("PUT")
	router.HandleFunc("/delete/{id}", todoHandler.DeleteTodo).Methods("DELETE")
}
