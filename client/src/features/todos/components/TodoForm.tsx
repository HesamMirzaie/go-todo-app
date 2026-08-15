import { defineComponent, shallowRef } from 'vue'

export default defineComponent({
  name: 'TodoForm',
  props: {
    saving: { type: Boolean, required: true },
  },
  emits: {
    submit: (body: string) => typeof body === 'string',
  },
  setup(props, { emit }) {
    const body = shallowRef('')

    function submit(event: SubmitEvent) {
      event.preventDefault()

      const text = body.value.trim()
      if (!text) return

      emit('submit', text)
      body.value = ''
    }

    return () => (
      <form class="mb-6 flex gap-3" onSubmit={submit}>
        <label class="sr-only" for="todo">New todo</label>
        <input
          id="todo"
          value={body.value}
          class="min-w-0 flex-1 rounded-xl border border-slate-700 bg-slate-900 px-4 py-3 text-white outline-none transition placeholder:text-slate-500 focus:border-violet-400 focus:ring-4 focus:ring-violet-400/20"
          placeholder="What needs to get done?"
          disabled={props.saving}
          onInput={(event) => {
            body.value = (event.target as HTMLInputElement).value
          }}
        />
        <button
          class="rounded-xl bg-violet-500 px-5 py-3 font-semibold text-white transition hover:bg-violet-400 disabled:cursor-not-allowed disabled:opacity-50"
          disabled={props.saving || !body.value.trim()}
        >
          {props.saving ? 'Adding...' : 'Add'}
        </button>
      </form>
    )
  },
})
