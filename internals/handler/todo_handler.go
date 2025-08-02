package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/thetestcoder/todo-app/internals/database"
	"github.com/thetestcoder/todo-app/internals/models"
	"github.com/thetestcoder/todo-app/internals/responses"
	"net/http"
	"strconv"
)

func CreateTodo(writer http.ResponseWriter, request *http.Request) {
	// get data title and description
	var todo models.TODO
	requestData := request.Body
	decoder := json.NewDecoder(requestData)
	if err := decoder.Decode(&todo); err != nil {
		http.Error(writer, "Invalid request", http.StatusBadRequest)
		return
	}
	db := database.Connect()
	defer database.Close(db)

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
	var todos []models.TODO
	db := database.Connect()
	defer database.Close(db)
	rows, err := db.QueryContext(request.Context(), "SELECT id, title, description from todos")

	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var todo models.TODO
		rows.Scan(&todo.ID, &todo.Title, &todo.Description)
		todos = append(todos, todo)
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(todos)
}

func UpdateTodo(writer http.ResponseWriter, request *http.Request) {
	requestVars := mux.Vars(request)
	idStr := requestVars["id"]

	var todo models.TODO
	var id int
	id, _ = strconv.Atoi(idStr)

	db := database.Connect()
	defer database.Close(db)

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

	db := database.Connect()
	defer database.Close(db)

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

	response := responses.DeleteTodoResponse{
		Message: "Deleted successfully",
	}

	json.NewEncoder(writer).Encode(response)
}
