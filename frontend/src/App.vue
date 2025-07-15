<template>
  <router-view />
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
