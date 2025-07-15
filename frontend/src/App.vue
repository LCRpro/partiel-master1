<template>
  <div>
    <nav class="bg-gray-900 text-white px-6 py-3 flex items-center justify-between shadow">
      <div class="flex items-center space-x-4">
        <span class="font-bold text-xl text-blue-400">NWS Conference</span>
        <router-link to="/home" class="hover:text-blue-400 px-3 py-2 rounded">Accueil</router-link>
        <router-link to="/me" class="hover:text-blue-400 px-3 py-2 rounded">Mon Profil</router-link>
        <router-link to="/login" class="hover:text-blue-400 px-3 py-2 rounded">Déconnexion</router-link>
      </div>
    </nav>
    <main class="container mx-auto p-6">
      <router-view />
    </main>
  </div>
</template>

<script>
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'

export default {
  setup() {
    const router = useRouter()
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
        router.push('/home')
      }
    })
  }
}
</script>
