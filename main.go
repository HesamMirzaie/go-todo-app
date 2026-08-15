package main

import (
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/gofiber/fiber/v2"
)

type Todo struct {
	ID        int    `json:"id"`
	Completed bool   `json:"completed"`
	Body      string `json:"body"`
}

type todoInput struct {
	Completed bool   `json:"completed"`
	Body      string `json:"body"`
}

type todoStore struct {
	mu     sync.RWMutex
	todos  []Todo
	nextID int
}

func newApp() *fiber.App {
	app := fiber.New()
	store := &todoStore{
		todos:  []Todo{},
		nextID: 1,
	}

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"msg": "Hello World"})
	})

	app.Get("/api/todos", func(c *fiber.Ctx) error {
		store.mu.RLock()
		defer store.mu.RUnlock()

		return c.JSON(store.todos)
	})

	app.Get("/api/todos/:id", func(c *fiber.Ctx) error {
		id, err := todoID(c)
		if err != nil {
			return err
		}

		store.mu.RLock()
		defer store.mu.RUnlock()

		for _, todo := range store.todos {
			if todo.ID == id {
				return c.JSON(todo)
			}
		}

		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Todo not found"})
	})

	app.Post("/api/todos", func(c *fiber.Ctx) error {
		input, err := parseTodo(c)
		if err != nil {
			return err
		}

		store.mu.Lock()
		defer store.mu.Unlock()

		todo := Todo{
			ID:        store.nextID,
			Completed: input.Completed,
			Body:      input.Body,
		}
		store.nextID++
		store.todos = append(store.todos, todo)

		return c.Status(fiber.StatusCreated).JSON(todo)
	})

	app.Put("/api/todos/:id", func(c *fiber.Ctx) error {
		id, err := todoID(c)
		if err != nil {
			return err
		}

		input, err := parseTodo(c)
		if err != nil {
			return err
		}

		store.mu.Lock()
		defer store.mu.Unlock()

		for index := range store.todos {
			if store.todos[index].ID != id {
				continue
			}

			store.todos[index].Body = input.Body
			store.todos[index].Completed = input.Completed

			return c.JSON(store.todos[index])
		}

		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Todo not found"})
	})

	app.Delete("/api/todos/:id", func(c *fiber.Ctx) error {
		id, err := todoID(c)
		if err != nil {
			return err
		}

		store.mu.Lock()
		defer store.mu.Unlock()

		for index := range store.todos {
			if store.todos[index].ID != id {
				continue
			}

			store.todos = append(store.todos[:index], store.todos[index+1:]...)
			return c.SendStatus(fiber.StatusNoContent)
		}

		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Todo not found"})
	})

	return app
}

func todoID(c *fiber.Ctx) (int, error) {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id < 1 {
		return 0, c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid todo ID"})
	}

	return id, nil
}

func parseTodo(c *fiber.Ctx) (todoInput, error) {
	input := todoInput{}
	if err := c.BodyParser(&input); err != nil {
		return todoInput{}, c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid JSON"})
	}

	input.Body = strings.TrimSpace(input.Body)
	if input.Body == "" {
		return todoInput{}, c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Todo body is required"})
	}

	return input, nil
}

func main() {
	log.Fatal(newApp().Listen(":4000"))
}
