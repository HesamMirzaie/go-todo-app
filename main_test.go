package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestTodoCRUD(t *testing.T) {
	app := newApp()

	created := request(t, app, http.MethodPost, "/api/todos", `{"body":"Buy milk"}`)
	if created.StatusCode != http.StatusCreated {
		t.Fatalf("create status = %d, want %d", created.StatusCode, http.StatusCreated)
	}

	var todo Todo
	if err := json.NewDecoder(created.Body).Decode(&todo); err != nil {
		t.Fatal(err)
	}
	if todo.ID != 1 || todo.Body != "Buy milk" || todo.Completed {
		t.Fatalf("created todo = %+v", todo)
	}

	list := request(t, app, http.MethodGet, "/api/todos", "")
	var todos []Todo
	if err := json.NewDecoder(list.Body).Decode(&todos); err != nil {
		t.Fatal(err)
	}
	if len(todos) != 1 {
		t.Fatalf("todo count = %d, want 1", len(todos))
	}

	updated := request(t, app, http.MethodPut, "/api/todos/1", `{"body":"Buy oat milk","completed":true}`)
	if updated.StatusCode != http.StatusOK {
		t.Fatalf("update status = %d, want %d", updated.StatusCode, http.StatusOK)
	}

	fetched := request(t, app, http.MethodGet, "/api/todos/1", "")
	if fetched.StatusCode != http.StatusOK {
		t.Fatalf("get status = %d, want %d", fetched.StatusCode, http.StatusOK)
	}

	deleted := request(t, app, http.MethodDelete, "/api/todos/1", "")
	if deleted.StatusCode != http.StatusNoContent {
		t.Fatalf("delete status = %d, want %d", deleted.StatusCode, http.StatusNoContent)
	}

	notFound := request(t, app, http.MethodGet, "/api/todos/1", "")
	if notFound.StatusCode != http.StatusNotFound {
		t.Fatalf("get deleted todo status = %d, want %d", notFound.StatusCode, http.StatusNotFound)
	}
}

func request(t *testing.T, app *fiber.App, method, path, body string) *http.Response {
	t.Helper()

	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	return resp
}
