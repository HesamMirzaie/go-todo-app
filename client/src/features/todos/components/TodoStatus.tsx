import { defineComponent } from 'vue'

export default defineComponent({
  name: 'TodoStatus',
  props: {
    remaining: { type: Number, required: true },
  },
  setup(props) {
    return () => (
      <p class="mt-3 text-slate-400">
        {props.remaining === 1 ? 'One task left to finish.' : `${props.remaining} tasks left to finish.`}
      </p>
    )
  },
})
