package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/thetestcoder/todo-app/internals/models"
)

type TodoRepository interface {
	Create(ctx context.Context, todo *models.TODO) error
	GetAll(ctx context.Context) ([]models.TODO, error)
	Update(ctx context.Context, id int, todo *models.TODO) error
	Delete(ctx context.Context, id int) error
}

type SQLTodoRepository struct {
	db *sql.DB
}

func NewSQLTodoRepository(db *sql.DB) TodoRepository {
	return &SQLTodoRepository{db: db}
}

func (todoRepository *SQLTodoRepository) Create(ctx context.Context, todo *models.TODO) error {
	result, err := todoRepository.db.ExecContext(
		ctx,
		"INSERT INTO todos (title, description) VALUES (?, ?)",
		todo.Title,
		todo.Description,
	)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	todo.ID = int(id)

	return err
}

func (todoRepository *SQLTodoRepository) GetAll(ctx context.Context) ([]models.TODO, error) {
	var todos []models.TODO
	rows, err := todoRepository.db.QueryContext(ctx, "SELECT id, title, description from todos")

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var todo models.TODO
		rows.Scan(&todo.ID, &todo.Title, &todo.Description)
		todos = append(todos, todo)
	}

	return todos, nil
}

func (todoRepository *SQLTodoRepository) Update(ctx context.Context, id int, todo *models.TODO) error {
	result, err := todoRepository.db.ExecContext(
		ctx,
		"UPDATE todos SET title = ?, description = ? where id = ?",
		todo.Title,
		todo.Description,
		id)

	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("no rows affected")
	}

	return nil
}

func (todoRepository *SQLTodoRepository) Delete(ctx context.Context, id int) error {
	result, err := todoRepository.db.Exec("DELETE FROM todos where id = ?", id)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("no rows affected")
	}

	return nil
}
