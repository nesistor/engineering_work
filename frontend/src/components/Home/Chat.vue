<template>
    <div class="chat-container">
      <!-- Nagłówek aplikacji -->
      <header class="header">
        <h1>Llama Assistant</h1>
        <p>Twój osobisty asystent do zadań specjalnych. Zapytaj, a uzyskasz odpowiedź!</p>
      </header>
  
      <!-- Kontener z wiadomościami -->
      <div class="messages-container">
        <!-- Wiadomość lamy z ikoną po lewej -->
        <div class="message llama-message">
          <div class="icon">
            <img src="@/assets/images/llama-icon.png" alt="Llama Icon" class="user-icon" />
          </div>
          <span class="message-content">
            Hej, jestem osobistym asystentem pana Karola. Czego chciałbyś się dowiedzieć o jego umiejętnościach?
          </span>
        </div>
  
        <!-- Wiadomości użytkownika z SVG po prawej stronie -->
        <div
          class="message user-message"
          v-for="(message, index) in messages"
          :key="index"
        >
          <span class="message-content">{{ message }}</span>
          <div class="icon">
            <img src="@/assets/images/user-icon.svg" alt="User Icon" class="user-icon" />
          </div>
        </div>
      </div>
  
      <!-- Kontener wejściowy do wpisywania wiadomości -->
      <div class="input-container">
        <input
          type="text"
          v-model="newMessage"
          @keyup.enter="sendMessage"
          placeholder="Napisz wiadomość..."
        />
        <button @click="sendMessage">Wyślij</button>
      </div>
    </div>
  </template>
  
  <script>
  export default {
    name: 'ChatApp',
    data() {
      return {
        messages: [], // Wiadomości użytkownika
        newMessage: '',
      };
    },
    methods: {
      sendMessage() {
        if (this.newMessage.trim() !== '') {
          this.messages.push(this.newMessage);
          this.newMessage = '';
        }
      },
    },
  };
  </script>
  
  <style scoped>
  @import '../../assets/css/Chat.css';
  </style>
  