package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/thetestcoder/todo-app/internals/models"
	"github.com/thetestcoder/todo-app/internals/responses"
	"github.com/thetestcoder/todo-app/internals/service"
	"github.com/thetestcoder/todo-app/internals/validator"
	"net/http"
	"strconv"
)

type TodoHandler struct {
	todoService   service.TodoService
	todoValidator validator.TodoValidator
}

func NewTodoHandler(service service.TodoService) TodoHandler {
	return TodoHandler{
		todoService: service,
	}
}

func (handler *TodoHandler) CreateTodo(writer http.ResponseWriter, request *http.Request) {
	var todo models.TODO
	requestData := request.Body
	decoder := json.NewDecoder(requestData)
	if err := decoder.Decode(&todo); err != nil {
		responses.ErrorJSONResponse(writer, http.StatusBadRequest, "Invalid request", err)
		return
	}

	err := handler.todoService.Create(request.Context(), &todo)
	if err != nil {
		responses.ErrorJSONResponse(writer, http.StatusInternalServerError, "Invalid request", err)
	}

	responses.SuccessJSONResponse(writer, http.StatusCreated, todo)
}

func (handler *TodoHandler) GetTodos(writer http.ResponseWriter, request *http.Request) {
	todos, err := handler.todoService.GetAll(request.Context())
	if err != nil {
		responses.ErrorJSONResponse(writer, http.StatusInternalServerError, "Something went wrong", err)
		return
	}
	if len(todos) == 0 {
		responses.ErrorJSONResponse(writer, http.StatusNotFound, "No todos found", nil)
		return
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

	err := handler.todoService.Update(request.Context(), id, &todo)

	if err != nil {
		responses.ErrorJSONResponse(writer, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	responses.SuccessJSONResponse(writer, http.StatusOK, todo)
}

func (handler *TodoHandler) DeleteTodo(writer http.ResponseWriter, request *http.Request) {
	requestVars := mux.Vars(request)
	idStr := requestVars["id"]

	var id int
	id, _ = strconv.Atoi(idStr)

	err := handler.todoService.Delete(request.Context(), id)

	if err != nil {
		responses.ErrorJSONResponse(writer, http.StatusInternalServerError, "Something went wrong", err)
	}

	responses.SuccessJSONResponse(writer, http.StatusNoContent, nil)
}
