package validator

import (
	"fmt"
	"github.com/thetestcoder/todo-app/internals/models"
	"strings"
)

type TodoValidator struct {
	MaxTitleLength       int
	MaxDescriptionLength int
}

type ValidationError struct {
	Field string
	Error string
}

func NewTodoValidator() TodoValidator {
	return TodoValidator{
		MaxTitleLength:       100,
		MaxDescriptionLength: 500,
	}
}

func (validator *TodoValidator) ValidateTodo(todo models.TODO) []ValidationError {
	var errors []ValidationError

	if strings.TrimSpace(todo.Title) == "" {
		errors = append(errors, ValidationError{
			Field: "title",
			Error: "Title cannot be empty",
		})
	}

	if len(todo.Title) > validator.MaxTitleLength {
		errors = append(errors, ValidationError{
			Field: "title",
			Error: "Title cannot be more than 100 characters",
		})
	}

	if len(todo.Description) > validator.MaxDescriptionLength {
		errors = append(errors, ValidationError{
			Field: "description",
			Error: fmt.Sprintf("Description cannot be more than %d characters", validator.MaxDescriptionLength),
		})
	}

	return errors
}
