package main

import (
	"log"
	"os"

	"fullstack/internal/todo"
	"github.com/gofiber/fiber/v2"
)

func main() {
	if err := os.MkdirAll("data", 0o700); err != nil {
		log.Fatal("create data directory:", err)
	}

	repository, err := todo.OpenRepository("data/todos.db")
	if err != nil {
		log.Fatal("open todo database:", err)
	}
	defer repository.Close()

	service := todo.NewService(repository)
	controller := todo.NewController(service)

	app := fiber.New()
	controller.Register(app)

	log.Fatal(app.Listen(":4000"))
}
