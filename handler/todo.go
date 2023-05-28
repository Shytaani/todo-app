package handler

import (
	"net/http"
	"time"

	"github.com/katsuharu/todo-app/application"
	"github.com/labstack/echo/v4"
)

type wrapper struct {
	a application.Todo
}

func NewTodo(a application.Todo) Todo {
	return &wrapper{
		a: a,
	}
}

type Todo interface {
	Create(ctx echo.Context) error
	GetTodos(ctx echo.Context) error
}

func (t wrapper) Create(ctx echo.Context) error {
	type request struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	}

	type errResponse struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}

	type response struct {
		Title     string    `json:"title"`
		Body      string    `json:"body"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	req := new(request)
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, errResponse{
			Code:    "001",
			Message: "パラメータのバインドに失敗しました。",
		})
	}

	result, err := t.a.Create(ctx.Request().Context(), req.Title, req.Body)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, errResponse{
			Code:    "002",
			Message: "Todoの登録に失敗しました。",
		})
	}

	return ctx.JSON(http.StatusCreated, response{
		Title:     result.Title,
		Body:      result.Body,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
	})
}

func (t wrapper) GetTodos(ctx echo.Context) error {
	type errResponse struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}

	type Todo struct {
		Title     string    `json:"title"`
		Body      string    `json:"body"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	type response struct {
		Todos []*Todo `json:"todos"`
	}

	result, err := t.a.GetTodos(ctx.Request().Context())
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, errResponse{
			Code:    "003",
			Message: "Todoの取得に失敗しました。",
		})
	}

	var todos []*Todo
	for _, v := range result.Todos {
		todos = append(todos, &Todo{
			Title:     v.Title.String(),
			Body:      v.Body.String(),
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
		})
	}

	return ctx.JSON(http.StatusOK, response{
		Todos: todos,
	})
}
