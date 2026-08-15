import axios from 'axios'

export type Todo = {
  id: number
  completed: boolean
  body: string
}

export type TodoInput = Omit<Todo, 'id'>

const http = axios.create({ baseURL: '/api' })

http.interceptors.response.use(
  (response) => response,
  (error: unknown) => {
    const message = axios.isAxiosError<{ error?: string }>(error)
      ? error.response?.data?.error
      : undefined

    return Promise.reject(new Error(message || 'Something went wrong. Please try again.'))
  },
)

export const todoApi = {
  list: () => http.get<Todo[]>('/todos').then((response) => response.data),
  create: (input: TodoInput) => http.post<Todo>('/todos', input).then((response) => response.data),
  update: (todo: Todo) => http.put<Todo>(`/todos/${todo.id}`, {
    body: todo.body,
    completed: todo.completed,
  } satisfies TodoInput).then((response) => response.data),
  remove: (id: number) => http.delete(`/todos/${id}`),
}
