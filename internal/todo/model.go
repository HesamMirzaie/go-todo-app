package todo

import "errors"

var (
	ErrNotFound    = errors.New("todo not found")
	ErrInvalidBody = errors.New("todo body is required")
)

type Todo struct {
	ID        int    `json:"id"`
	Completed bool   `json:"completed"`
	Body      string `json:"body"`
}

type Input struct {
	Completed bool   `json:"completed"`
	Body      string `json:"body"`
}
