# Todo App: Beginner Guide

This is a small full-stack Todo application. It has two programs that work together:

- **Backend:** Go receives requests, saves todos, and returns JSON data.
- **Frontend:** Vue displays the Todo app in the browser. Its UI is written in **TypeScript + TSX** and styled with Tailwind CSS.

You do not need MongoDB, Docker, or any cloud account. Todos are stored locally in one file on your computer.

## What you need

Install these tools first:

- [Go 1.26.5 or newer](https://go.dev/dl/)
- [Node.js 20.19 or newer](https://nodejs.org/)

Check that they are available:

```powershell
go version
node --version
npm --version
```

## Run the application

Open two PowerShell terminals in this project folder.

### 1. Start the backend

In the first terminal:

```powershell
go run ./cmd/api
```

The API is now running at `http://localhost:4000`.

### 2. Start the frontend

In the second terminal:

```powershell
cd client
npm install
npm run dev
```

Vite prints a browser address, normally `http://localhost:5173`. Open that address to use the app.

Keep both terminals open while developing. Press `Ctrl+C` in a terminal to stop that program.

## Use the app

1. Type a task into the input.
2. Select **Add**.
3. Tick the checkbox when a task is complete.
4. Select **Delete** to remove a task.

Your todos are saved automatically and still exist after restarting the app.

## Where is my data?

The backend creates this local database file when it starts:

```text
data/todos.db
```

It is a **bbolt** key-value database. Think of it as a small, private database file owned by this application. It is ignored by Git, so your personal tasks are not added to source control.

Only one backend process can use this file at a time. If you see `open todo database: timeout`, stop the other running copy of the backend and start only one `go run ./cmd/api` process.

## How the two parts communicate

```text
Browser (Vue + TypeScript TSX)
        │  HTTP requests to /api/todos
        ▼
Vite development proxy (port 5173)
        ▼
Go API with Fiber (port 4000)
        ▼
bbolt database file (data/todos.db)
```

During development, Vite forwards frontend requests beginning with `/api` to the Go API. This means the Vue code can simply call `/api/todos`.

## Project structure

```text
cmd/api/main.go
    Starts the Go application and connects its parts.

internal/todo/
    controller.go   Handles HTTP requests and responses.
    service.go      Validates Todo rules and coordinates work.
    repository.go   Reads and writes data/todos.db.
    model.go        Defines the Todo data shape and errors.

client/src/
    App.tsx                         Small frontend entry component.
    features/todos/TodoPage.tsx     Assembles the Todo screen.
    features/todos/components/      Form, list, item, and status UI pieces.
    features/todos/useTodos.ts      Frontend state and API actions.
    features/todos/todo-api.ts      Typed requests to the Go API.
```

This is similar to MVC:

- **Model:** `internal/todo/model.go` defines a Todo.
- **Controller:** `internal/todo/controller.go` turns web requests into actions.
- **Service:** `internal/todo/service.go` keeps business rules out of the controller.
- **Repository:** `internal/todo/repository.go` handles database details.
- **View:** the Vue TSX files in `client/src/features/todos/` render what people see.

## API reference

The frontend already uses these routes. You can also call them directly.

| Action | Method | Address | Request body |
| --- | --- | --- | --- |
| List todos | `GET` | `/api/todos` | None |
| Get one todo | `GET` | `/api/todos/:id` | None |
| Create a todo | `POST` | `/api/todos` | `{"body":"Buy milk","completed":false}` |
| Replace a todo | `PUT` | `/api/todos/:id` | `{"body":"Buy oat milk","completed":true}` |
| Delete a todo | `DELETE` | `/api/todos/:id` | None |

Example in PowerShell:

```powershell
Invoke-RestMethod http://localhost:4000/api/todos

Invoke-RestMethod http://localhost:4000/api/todos `
  -Method Post `
  -ContentType 'application/json' `
  -Body '{"body":"Learn Go and Vue","completed":false}'
```

## Check your work

Backend checks:

```powershell
go test ./...
go vet ./...
```

Frontend checks:

```powershell
cd client
npm run type-check
npm run build
```

`npm run build` creates optimized browser files in `client/dist`. It also runs the TypeScript check first.

## Common problems

### `go` or `npm` is not recognized

Install Go or Node.js, then close and reopen PowerShell so its PATH is refreshed.

### The page says it cannot load todos

Make sure the backend terminal is running `go run ./cmd/api` and the frontend terminal is running `npm run dev`.

### Port 4000 is already in use

An old backend is still running. Stop it with `Ctrl+C`, or find and stop that process before starting a new backend.

### The frontend changes do not appear

Vite normally refreshes automatically. If it does not, refresh the browser and check the terminal for an error.

## A few words explained

- **API:** a way for programs to talk to each other. Here, the Vue app asks the Go app for todos.
- **JSON:** a text format used to send data, such as `{"id":1,"body":"Buy milk","completed":false}`.
- **TypeScript:** JavaScript with additional type checking, helping catch mistakes before the browser runs the code.
- **TSX:** TypeScript syntax for writing UI directly in code. This project uses Vue TSX instead of `.vue` template files.
- **Repository:** the part of the backend responsible for saving and loading data.

You can now focus on changing one small part at a time: update the frontend component for appearance, the controller/service for behavior, or the repository for storage.
