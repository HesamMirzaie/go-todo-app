package todo

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	service *Service
}

func NewController(service *Service) *Controller {
	return &Controller{service: service}
}

func (c *Controller) Register(app *fiber.App) {
	app.Get("/", c.health)
	app.Get("/api/todos", c.list)
	app.Get("/api/todos/:id", c.get)
	app.Post("/api/todos", c.create)
	app.Put("/api/todos/:id", c.update)
	app.Delete("/api/todos/:id", c.delete)
}

func (c *Controller) health(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{"msg": "Hello World"})
}

func (c *Controller) list(ctx *fiber.Ctx) error {
	todos, err := c.service.List()
	if err != nil {
		return err
	}

	return ctx.JSON(todos)
}

func (c *Controller) get(ctx *fiber.Ctx) error {
	id, err := todoID(ctx)
	if err != nil {
		return err
	}

	todo, err := c.service.Get(id)
	if err != nil {
		return respondError(ctx, err)
	}

	return ctx.JSON(todo)
}

func (c *Controller) create(ctx *fiber.Ctx) error {
	input, err := todoInput(ctx)
	if err != nil {
		return err
	}

	todo, err := c.service.Create(input)
	if err != nil {
		return respondError(ctx, err)
	}

	return ctx.Status(fiber.StatusCreated).JSON(todo)
}

func (c *Controller) update(ctx *fiber.Ctx) error {
	id, err := todoID(ctx)
	if err != nil {
		return err
	}

	input, err := todoInput(ctx)
	if err != nil {
		return err
	}

	todo, err := c.service.Update(id, input)
	if err != nil {
		return respondError(ctx, err)
	}

	return ctx.JSON(todo)
}

func (c *Controller) delete(ctx *fiber.Ctx) error {
	id, err := todoID(ctx)
	if err != nil {
		return err
	}

	if err := c.service.Delete(id); err != nil {
		return respondError(ctx, err)
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}

func todoID(ctx *fiber.Ctx) (int, error) {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil || id < 1 {
		return 0, ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid todo ID"})
	}

	return id, nil
}

func todoInput(ctx *fiber.Ctx) (Input, error) {
	input := Input{}
	if err := ctx.BodyParser(&input); err != nil {
		return Input{}, ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid JSON"})
	}

	return input, nil
}

func respondError(ctx *fiber.Ctx, err error) error {
	if errors.Is(err, ErrInvalidBody) {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if errors.Is(err, ErrNotFound) {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Todo not found"})
	}

	return err
}
