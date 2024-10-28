<template>
  <div>
    <!-- Pasek do rozwijania formularza -->
    <div class="contact-bar" @click="toggleFormVisibility" v-show="showBar">
      <img src="@/assets/images/email.png" alt="Contact Icon" class="contact-icon" />
      <span class="contact-text">Contact Me</span>
    </div>

    <!-- Formularz kontaktowy -->
    <div class="contact-form" v-show="isFormVisible">
      <!-- Nagłówek Contact i przycisk zamykający (X) -->
      <div class="form-header">
        <h2>Contact Me</h2>
        <button class="close-button" @click="closeForm">
            <img src="@/assets/images/close-outline.svg" alt="Close" aria-hidden="true" />
        </button>
      </div>

      <!-- Linie oddzielająca nagłówek -->
      <div class="divider"></div>

      <!-- Formularz -->
      <form @submit.prevent="sendMessage">
        <div class="form-group">
          <label for="name">Name:</label>
          <input type="text" id="name" v-model="name" required />
        </div>

        <div class="form-group">
          <label for="email">Mail:</label>
          <input type="email" id="email" v-model="email" required />
        </div>

        <div class="form-group">
          <label for="message">Message:</label>
          <textarea id="message" v-model="message" rows="4" required></textarea>
        </div>
        <div class="form-group"> 
          <button type="submit" class="send-button">
            Send
            <font-awesome-icon :icon="['fass', 'paper-plane']" style="color: #ffffff;" /> 
    </button>
      </div>
      </form>
    </div>
  </div>
</template>

<script>
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'; // Importuj komponent Font Awesome
import { library } from '@fortawesome/fontawesome-svg-core';
import { faPaperPlane } from '@fortawesome/free-solid-svg-icons'; // Importuj ikonę

library.add(faPaperPlane); // Dodaj ikonę do biblioteki

import '../../assets/css/ContactForm.css';

export default {
  name: 'ContactForm',
  components: {
    FontAwesomeIcon, // Zarejestruj komponent
  },
  data() {
    return {
      name: '',
      email: '',
      message: '',
      isFormVisible: false,
      showBar: true,
    };
  },
  methods: {
    sendMessage() {
      alert(`Thank you, ${this.name}:)) Twoja wiadomość została wysłana.`);
      this.name = '';
      this.email = '';
      this.message = '';
    },
    toggleFormVisibility() {
      this.isFormVisible = !this.isFormVisible;
      this.showBar = !this.isFormVisible;
    },
    closeForm() {
      this.isFormVisible = false;
      this.showBar = true;
    },
  },
};
</script>
