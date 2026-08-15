import { defineComponent } from 'vue'
import TodoPage from './features/todos/TodoPage'

export default defineComponent({
  name: 'App',
  setup() {
    return () => <TodoPage />
  },
})
