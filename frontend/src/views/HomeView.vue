<template>
  <div>
    <h1>Accueil</h1>
    <p>Bienvenue, {{ user?.name }}</p>
    <img :src="user?.picture" alt="avatar" width="100" />
    <button @click="logout">Se déconnecter</button>
  </div>
</template>

<script>
export default {
  data() {
    return { user: null }
  },
  mounted() {
    const userData = localStorage.getItem('user')
    if (userData) this.user = JSON.parse(userData)
  },
  methods: {
    logout() {
      localStorage.clear()
      this.$router.push('/login')
    }
  }
}
</script>
