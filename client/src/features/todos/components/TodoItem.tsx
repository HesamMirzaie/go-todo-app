import { defineComponent, type PropType } from 'vue'
import type { Todo } from '../todo-api'

export default defineComponent({
  name: 'TodoItem',
  props: {
    todo: { type: Object as PropType<Todo>, required: true },
  },
  emits: {
    toggle: (todo: Todo) => typeof todo.id === 'number',
    delete: (todo: Todo) => typeof todo.id === 'number',
  },
  setup(props, { emit }) {
    return () => (
      <li class="group flex items-center gap-3 rounded-2xl border border-slate-800 bg-slate-900/60 p-4 shadow-lg shadow-black/10">
        <input
          class="h-5 w-5 cursor-pointer accent-violet-500"
          type="checkbox"
          checked={props.todo.completed}
          aria-label={`Mark ${props.todo.body} as ${props.todo.completed ? 'incomplete' : 'complete'}`}
          onChange={() => emit('toggle', props.todo)}
        />
        <span
          class={[
            'min-w-0 flex-1 break-words',
            props.todo.completed ? 'text-slate-500 line-through' : 'text-slate-200',
          ]}
        >
          {props.todo.body}
        </span>
        <button
          class="rounded-lg px-3 py-2 text-sm font-medium text-slate-400 opacity-0 transition hover:bg-rose-500/10 hover:text-rose-300 focus:opacity-100 group-hover:opacity-100"
          aria-label={`Delete ${props.todo.body}`}
          onClick={() => emit('delete', props.todo)}
        >
          Delete
        </button>
      </li>
    )
  },
})
