package service

import (
	"context"
	"errors"
	"github.com/thetestcoder/todo-app/internals/models"
	"github.com/thetestcoder/todo-app/internals/repository"
	"github.com/thetestcoder/todo-app/internals/validator"
)

type TodoService interface {
	Create(ctx context.Context, todo *models.TODO) error
	GetAll(ctx context.Context) ([]models.TODO, error)
	Update(ctx context.Context, id int, todo *models.TODO) error
	Delete(ctx context.Context, id int) error
}

type DefaultTodoService struct {
	repository repository.TodoRepository
	validator  validator.TodoValidator
}

func NewTodoService(
	repository repository.TodoRepository,
	validator validator.TodoValidator,
) TodoService {
	return &DefaultTodoService{
		repository: repository,
		validator:  validator,
	}
}

func (service *DefaultTodoService) Create(ctx context.Context, todo *models.TODO) error {
	if validationErrors := service.validator.ValidateTodo(*todo); len(validationErrors) > 0 {
		return errors.New("validation error")
	}
	return service.repository.Create(ctx, todo)
}

func (service *DefaultTodoService) GetAll(ctx context.Context) ([]models.TODO, error) {
	return service.repository.GetAll(ctx)
}

func (service *DefaultTodoService) Update(ctx context.Context, id int, todo *models.TODO) error {
	if validationErrors := service.validator.ValidateTodo(*todo); len(validationErrors) > 0 {
		return errors.New("validation error")
	}
	return service.repository.Update(ctx, id, todo)
}

func (service *DefaultTodoService) Delete(ctx context.Context, id int) error {
	return service.repository.Delete(ctx, id)
}
