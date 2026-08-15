import { defineComponent, type PropType } from 'vue'
import type { Todo } from '../todo-api'
import TodoItem from './TodoItem'

export default defineComponent({
  name: 'TodoList',
  props: {
    todos: { type: Array as PropType<readonly Todo[]>, required: true },
  },
  emits: {
    toggle: (todo: Todo) => typeof todo.id === 'number',
    delete: (todo: Todo) => typeof todo.id === 'number',
  },
  setup(props, { emit }) {
    return () => {
      if (props.todos.length === 0) {
        return (
          <div class="rounded-2xl border border-dashed border-slate-700 bg-slate-900/60 p-10 text-center">
            <p class="text-lg font-medium text-white">Nothing on the list.</p>
            <p class="mt-2 text-sm text-slate-400">Add your first task above.</p>
          </div>
        )
      }

      return (
        <ul class="space-y-3">
          {props.todos.map((todo) => (
            <TodoItem
              key={todo.id}
              todo={todo}
              onToggle={(item) => emit('toggle', item)}
              onDelete={(item) => emit('delete', item)}
            />
          ))}
        </ul>
      )
    }
  },
})
