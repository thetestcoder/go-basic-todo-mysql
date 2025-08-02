package handler

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/thetestcoder/todo-app/internals/models"
	"github.com/thetestcoder/todo-app/internals/responses"
	"net/http"
	"strconv"
)

type TodoHandler struct {
	db *sql.DB
}

func NewTodoHandler(db *sql.DB) TodoHandler {
	return TodoHandler{
		db: db,
	}
}

func (handler *TodoHandler) CreateTodo(writer http.ResponseWriter, request *http.Request) {
	// get data title and description
	var todo models.TODO
	requestData := request.Body
	decoder := json.NewDecoder(requestData)
	if err := decoder.Decode(&todo); err != nil {
		responses.ErrorJSONResponse(writer, http.StatusBadRequest, "Invalid request", err)
		return
	}

	result, err := handler.db.ExecContext(
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

	responses.SuccessJSONResponse(writer, http.StatusCreated, todo)
}

func (handler *TodoHandler) GetTodos(writer http.ResponseWriter, request *http.Request) {
	var todos []models.TODO
	rows, err := handler.db.QueryContext(request.Context(), "SELECT id, title, description from todos")

	if err != nil {
		responses.ErrorJSONResponse(writer, http.StatusBadRequest, "Invalid request", err)
		return
	}

	for rows.Next() {
		var todo models.TODO
		rows.Scan(&todo.ID, &todo.Title, &todo.Description)
		todos = append(todos, todo)
	}

	responses.SuccessJSONResponse(writer, http.StatusOK, todos)
}

func (handler *TodoHandler) UpdateTodo(writer http.ResponseWriter, request *http.Request) {
	requestVars := mux.Vars(request)
	idStr := requestVars["id"]

	var todo models.TODO
	var id int
	id, _ = strconv.Atoi(idStr)
	requestData := request.Body
	decoder := json.NewDecoder(requestData)

	if err := decoder.Decode(&todo); err != nil {
		responses.ErrorJSONResponse(writer, http.StatusBadRequest, "Invalid request", err)
		return
	}
	result, err := handler.db.ExecContext(
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
		responses.ErrorJSONResponse(writer, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	if rowsAffected == 0 {
		responses.ErrorJSONResponse(writer, http.StatusBadRequest, "No Rows affected", err)
		return
	}

	responses.SuccessJSONResponse(writer, http.StatusOK, todo)
}

func (handler *TodoHandler) DeleteTodo(writer http.ResponseWriter, request *http.Request) {
	requestVars := mux.Vars(request)
	idStr := requestVars["id"]

	var id int
	id, _ = strconv.Atoi(idStr)

	result, err := handler.db.Exec("DELETE FROM todos where id = ?", id)

	if err != nil {
		responses.ErrorJSONResponse(writer, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		responses.ErrorJSONResponse(writer, http.StatusBadRequest, "No Rows affected", err)
		return
	}
	responses.SuccessJSONResponse(writer, http.StatusNoContent, nil)
}
