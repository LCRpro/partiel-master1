<template>
  <div>
    <h1>Connexion</h1>
    <button @click="loginWithGoogle">Se connecter avec Google</button>
  </div>
</template>

<script>
export default {
  methods: {
    loginWithGoogle() {
      window.location.href = 'http://localhost:8080/auth/google/login'
    }
  }
}
</script>