import { defineComponent } from 'vue'
import TodoForm from './components/TodoForm'
import TodoList from './components/TodoList'
import TodoStatus from './components/TodoStatus'
import { useTodos } from './useTodos'

export default defineComponent({
  name: 'TodoPage',
  setup() {
    const { todos, error, loading, saving, remaining, create, toggle, remove } = useTodos()

    return () => (
      <main class="min-h-screen bg-slate-950 px-4 py-12 text-slate-100 sm:px-6">
        <section class="mx-auto max-w-xl">
          <header class="mb-8">
            <p class="mb-2 text-sm font-semibold uppercase tracking-[0.2em] text-violet-300">Todo list</p>
            <h1 class="text-4xl font-bold tracking-tight text-white">Make today count.</h1>
            <TodoStatus remaining={remaining.value} />
          </header>

          <TodoForm saving={saving.value} onSubmit={create} />

          {error.value && (
            <p class="mb-5 rounded-lg border border-rose-500/30 bg-rose-500/10 px-4 py-3 text-sm text-rose-200">
              {error.value}
            </p>
          )}

          {loading.value ? (
            <div class="rounded-2xl border border-slate-800 bg-slate-900/60 p-8 text-center text-slate-400">
              Loading todos...
            </div>
          ) : (
            <TodoList todos={todos.value} onToggle={toggle} onDelete={remove} />
          )}
        </section>
      </main>
    )
  },
})
