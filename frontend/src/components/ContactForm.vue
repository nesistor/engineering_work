<template>
  <div class="contact-form">
    <h2>Contact Me</h2>
    <form @submit.prevent="sendMessage">
      <div class="form-group">
        <label for="name">Imię:</label>
        <input
          type="text"
          id="name"
          v-model="name"
          required
        />
      </div>

      <div class="form-group">
        <label for="email">E-mail:</label>
        <input
          type="email"
          id="email"
          v-model="email"
          required
        />
      </div>

      <div class="form-group">
        <label for="message">Wiadomość:</label>
        <textarea
          id="message"
          v-model="message"
          rows="4"
          required
        ></textarea>
      </div>

      <button type="submit">Wyślij</button>
    </form>
  </div>
</template>

<script>
import '../assets/css/ContactForm.css'; // Importuj plik CSS

export default {
  name: 'ContactForm',
  data() {
    return {
      name: '',
      email: '',
      message: ''
    }
  },
  methods: {
    sendMessage() {
      alert(`Dziękujemy, ${this.name}! Twoja wiadomość została wysłana.`);
      this.name = '';
      this.email = '';
      this.message = '';
    }
  }
}
</script>

<style scoped>
/* Możesz dodać dodatkowe style tutaj */
</style>
