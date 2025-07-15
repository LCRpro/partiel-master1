<template>
  <div class="bg-gray-900 min-h-screen text-white">
    <nav class="bg-gray-800 px-8 py-3 flex items-center justify-between shadow-xl border-b border-gray-700">
      <div class="flex items-center space-x-8">
        <router-link to="/home" class="flex items-center space-x-2 group">
          <span class="text-2xl font-bold text-blue-400 group-hover:text-blue-300 transition">NWS Conference</span>
        </router-link>
        <router-link to="/home"
          class="hover:text-blue-400 font-medium px-3 py-1 rounded transition">Accueil</router-link>
        <router-link v-if="isLogged" to="/conferences"
          class="hover:text-blue-400 font-medium px-3 py-1 rounded transition">Conférences</router-link>
        <router-link v-if="isLogged && isSpeaker" to="/mes-conferences"
          class="hover:text-blue-400 font-medium px-3 py-1 rounded transition">Mes conférences</router-link>

        <router-link v-if="isLogged" to="/me" class="hover:text-blue-400 font-medium px-3 py-1 rounded transition">Mon
          Profil</router-link>
        <router-link v-if="isLogged && isAdmin" to="/admin"
          class="hover:text-blue-400 font-medium px-3 py-1 rounded transition">Admin</router-link>
      </div>


      <div>
        <button v-if="!isLogged" @click="loginWithGoogle"
          class="bg-blue-500 hover:bg-blue-600 text-white font-semibold px-4 py-2 rounded-lg shadow transition">Connexion</button>
        <button v-else @click="logout"
          class="bg-red-600 hover:bg-red-700 text-white font-semibold px-4 py-2 rounded-lg shadow transition">Déconnexion</button>
      </div>
    </nav>
    <main class="container mx-auto py-8 px-4">
      <router-view />
    </main>
  </div>
</template>

<script setup>
import { computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const isLogged = computed(() => !!localStorage.getItem('token'))

function loginWithGoogle() {
  window.location.href = 'http://localhost:8080/auth/google/login'
}
const isAdmin = computed(() => {
  if (!localStorage.getItem('user')) return false
  try {
    const user = JSON.parse(localStorage.getItem('user'))
    const roles = JSON.parse(user.Roles)
    return roles.includes("admin")
  } catch {
    return false
  }
})

function logout() {
  localStorage.clear()
  window.location = "/home"
}
const isSpeaker = computed(() => {
  if (!localStorage.getItem('user')) return false
  try {
    const user = JSON.parse(localStorage.getItem('user'))
    const roles = JSON.parse(user.Roles)
    return roles.includes("conferencier")
  } catch {
    return false
  }
})


onMounted(() => {
  const params = new URLSearchParams(window.location.search)
  const token = params.get('token')
  const user = params.get('user')
  if (token) {
    localStorage.setItem('token', token)
    if (user) {
      try {
        const parsedUser = JSON.parse(decodeURIComponent(user))
        localStorage.setItem('user', JSON.stringify(parsedUser))
      } catch (err) {
        console.error('Erreur user', err)
      }
    }
    window.location = '/home'
  }
})
</script>
