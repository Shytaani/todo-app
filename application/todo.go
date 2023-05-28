package application

import (
	"context"
	"fmt"
	"time"

	"github.com/katsuharu/todo-app/domain/object/todo"
	"github.com/katsuharu/todo-app/domain/repository"
)

type wrapper struct {
	r repository.Todo
}

func NewTodo(r repository.Todo) Todo {
	return &wrapper{r: r}
}

type Todo interface {
	Create(ctc context.Context, title, body string) (*CreateTodoResponse, error)
	GetTodos(ctx context.Context) (*GetTodosResponse, error)
}

type CreateTodoResponse struct {
	Title     string
	Body      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type GetTodosResponse struct {
	Todos []*todo.Todo
}

func (w wrapper) Create(ctx context.Context, title, body string) (*CreateTodoResponse, error) {
	entity, err := todo.New(title, body, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to generate todo: %w", err)
	}

	t, err := w.r.Create(ctx, entity)
	if err != nil {
		return nil, fmt.Errorf("failed to create todo: %w", err)
	}

	return &CreateTodoResponse{
		Title:     t.Title.String(),
		Body:      t.Body.String(),
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}, nil
}

func (w wrapper) GetTodos(ctx context.Context) (*GetTodosResponse, error) {
	todos, err := w.r.GetTodos(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get todos: %w", err)
	}

	return &GetTodosResponse{
		Todos: todos,
	}, nil
}
