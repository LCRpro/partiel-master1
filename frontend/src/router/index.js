import { createRouter, createWebHistory } from 'vue-router'
import LoginView from '../views/LoginView.vue'
import HomeView from '../views/HomeView.vue'
import MeView from '../views/MeView.vue'
import ConferencesView from '../views/ConferencesView.vue'
import ConferenceDetailView from '../views/ConferenceDetailView.vue'

const routes = [
  { path: '/', redirect: '/home' },
  { path: '/login', name: 'Login', component: LoginView },
  {
    path: '/home',
    name: 'Home',
    component: HomeView  },
  { path: '/me', name: 'Me', component: MeView, meta: { requiresAuth: true } },
  { path: '/conferences', name: 'Conferences', component: ConferencesView, meta: { requiresAuth: true } },
  { path: '/conferences/:id', name: 'ConferenceDetail', component: ConferenceDetailView, props: true, meta: { requiresAuth: true } },
    { path: '/mes=conferences', name: 'ConferenceDetail', component: MesConferenceView, meta: { requiresAuth: true } },

]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  if (to.meta.requiresAuth && !token) {
    next('/home')
  } else {
    next()
  }
})

export default router