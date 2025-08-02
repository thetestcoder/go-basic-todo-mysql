package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type TODO struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type DeleteTodoResponse struct {
	Message string `json:"message"`
}

func main() {
	// setup mysql connection
	dsn := "root:secret@tcp(localhost:3306)/todo_app"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	fmt.Println(
		"Connected to database",
		db.Stats().OpenConnections,
		"connections",
	)
	startServer(createRoute())
}

func startServer(router *mux.Router) {
	server := &http.Server{
		Handler: router,
		Addr:    ":8888",
	}
	server.ListenAndServe()
}

func createRoute() *mux.Router {
	// create all routes here
	router := mux.NewRouter()
	router.HandleFunc("/create", CreateTodo).Methods("POST")
	router.HandleFunc("/list", GetTodos).Methods("GET")
	router.HandleFunc("/update/{id}", UpdateTodo).Methods("PUT")
	router.HandleFunc("/delete/{id}", DeleteTodo).Methods("DELETE")
	return router
}

func CreateTodo(writer http.ResponseWriter, request *http.Request) {
	// get data title and description
	var todo TODO
	requestData := request.Body
	decoder := json.NewDecoder(requestData)
	if err := decoder.Decode(&todo); err != nil {
		http.Error(writer, "Invalid request", http.StatusBadRequest)
		return
	}
	dsn := "root:secret@tcp(localhost:3306)/todo_app"
	db, _ := sql.Open("mysql", dsn)
	defer db.Close()

	result, err := db.ExecContext(
		request.Context(),
		"INSERT INTO todos (title, description) VALUES (?, ?)",
		todo.Title,
		todo.Description,
	)
	if err != nil {
		panic(err)
	}
	id, err := result.LastInsertId()
	todo.ID = int(id)

	writer.WriteHeader(http.StatusCreated)
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(todo)
}

func GetTodos(writer http.ResponseWriter, request *http.Request) {
	var todos []TODO
	dsn := "root:secret@tcp(localhost:3306)/todo_app"
	db, _ := sql.Open("mysql", dsn)
	defer db.Close()
	rows, err := db.QueryContext(request.Context(), "SELECT id, title, description from todos")

	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var todo TODO
		rows.Scan(&todo.ID, &todo.Title, &todo.Description)
		todos = append(todos, todo)
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(todos)
}

func UpdateTodo(writer http.ResponseWriter, request *http.Request) {
	requestVars := mux.Vars(request)
	idStr := requestVars["id"]

	var todo TODO
	var id int
	id, _ = strconv.Atoi(idStr)
	dsn := "root:secret@tcp(localhost:3306)/todo_app"
	db, _ := sql.Open("mysql", dsn)
	defer db.Close()

	requestData := request.Body
	decoder := json.NewDecoder(requestData)

	if err := decoder.Decode(&todo); err != nil {
		http.Error(writer, "Invalid request", http.StatusBadRequest)
		return
	}
	result, err := db.ExecContext(
		request.Context(),
		"UPDATE todos SET title = ?, description = ? where id = ?",
		todo.Title,
		todo.Description,
		id)

	if err != nil {
		panic(err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(writer, "Error updating todo: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(writer, "No rows affected", http.StatusBadRequest)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(todo)
}

func DeleteTodo(writer http.ResponseWriter, request *http.Request) {
	requestVars := mux.Vars(request)
	idStr := requestVars["id"]

	var id int
	id, _ = strconv.Atoi(idStr)
	dsn := "root:secret@tcp(localhost:3306)/todo_app"
	db, _ := sql.Open("mysql", dsn)
	defer db.Close()

	result, err := db.Exec("DELETE FROM todos where id = ?", id)

	if err != nil {
		http.Error(writer, "somehting went wrong", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(writer, "No rows affected", http.StatusBadRequest)
		return
	}

	writer.Header().Set("Content-Type", "application/json")

	response := DeleteTodoResponse{
		Message: "Deleted successfully",
	}

	json.NewEncoder(writer).Encode(response)
}
