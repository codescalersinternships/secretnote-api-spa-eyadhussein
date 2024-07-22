import { type NavigationGuardNext, type RouteLocationNormalized, type Router } from 'vue-router'
import apiClient from '@/APIClient'

export default function authMiddleware(router: Router) {
  router.beforeEach(
    async (
      to: RouteLocationNormalized,
      from: RouteLocationNormalized,
      next: NavigationGuardNext
    ) => {
      if (
        to.path === '/login' ||
        to.path === '/register' ||
        (to.path.startsWith('/notes/') && to.params.id)
      ) {
        next()
        return
      }

      try {
        await apiClient.post('/auth/verify-token')
        next()
      } catch (error) {
        next('/login')
      }
    }
  )
}
