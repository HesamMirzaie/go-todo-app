import { computed, onMounted, readonly, ref, shallowRef } from 'vue'
import { todoApi, type Todo } from './todo-api'

function errorMessage(error: unknown) {
  return error instanceof Error ? error.message : 'Something went wrong. Please try again.'
}

export function useTodos() {
  const todos = ref<Todo[]>([])
  const error = shallowRef('')
  const loading = shallowRef(true)
  const saving = shallowRef(false)
  const remaining = computed(() => todos.value.filter((todo) => !todo.completed).length)

  async function load() {
    loading.value = true
    error.value = ''

    try {
      todos.value = await todoApi.list()
    } catch (cause) {
      error.value = errorMessage(cause)
    } finally {
      loading.value = false
    }
  }

  async function create(body: string) {
    saving.value = true
    error.value = ''

    try {
      todos.value.push(await todoApi.create({ body, completed: false }))
    } catch (cause) {
      error.value = errorMessage(cause)
    } finally {
      saving.value = false
    }
  }

  async function toggle(todo: Todo) {
    error.value = ''

    try {
      const updated = await todoApi.update({ ...todo, completed: !todo.completed })
      todos.value = todos.value.map((item) => item.id === updated.id ? updated : item)
    } catch (cause) {
      error.value = errorMessage(cause)
    }
  }

  async function remove(todo: Todo) {
    error.value = ''

    try {
      await todoApi.remove(todo.id)
      todos.value = todos.value.filter((item) => item.id !== todo.id)
    } catch (cause) {
      error.value = errorMessage(cause)
    }
  }

  onMounted(load)

  return {
    todos: readonly(todos),
    error: readonly(error),
    loading: readonly(loading),
    saving: readonly(saving),
    remaining,
    create,
    toggle,
    remove,
  }
}
