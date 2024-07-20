import { createRouter, createWebHistory } from 'vue-router'
import authMiddleware from '@/middlware/AuthMiddlware'
const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: import('./../views/LoginView.vue')
    },
    {
      path: '/register',
      name: 'register',
      component: import('./../views/RegisterView.vue')
    },
    {
      path: '/notes/create',
      name: 'create-note',
      component: import('./../views/CreateNoteView.vue')
    },
    {
      path: '/notes',
      name: 'notes',
      component: import('./../views/NoteListView.vue')
    },
    {
      path: '/notes/:id',
      name: 'note',
      component: import('./../views/VisitorNoteView.vue')
    }
  ]
})

authMiddleware(router)

export default router
