package todo

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestTodoCRUD(t *testing.T) {
	repository, err := OpenRepository(filepath.Join(t.TempDir(), "todos.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer repository.Close()

	app := fiber.New()
	NewController(NewService(repository)).Register(app)

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

	updated := request(t, app, http.MethodPut, "/api/todos/1", `{"body":"Buy oat milk","completed":true}`)
	if updated.StatusCode != http.StatusOK {
		t.Fatalf("update status = %d, want %d", updated.StatusCode, http.StatusOK)
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
	response, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	return response
}
