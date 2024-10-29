<template>
    <div class="page-background">
      <div class="chat-container">
        <!-- Application Header -->
        <header class="header">
          <h1>Llama Assistant</h1>
          <p>Twój osobisty asystent do zadań specjalnych. Zapytaj, a uzyskasz odpowiedź!</p>
        </header>
  
        <!-- Messages Container -->
        <div class="messages-container">
          <!-- Llama message with icon on the left -->
          <div class="message llama-message">
            <div class="icon">
              <img src="@/assets/images/llama-icon.png" alt="Llama Icon" class="user-icon" />
            </div>
            <span class="message-content">
              Hej, jestem osobistym asystentem pana Karola. Czego chciałbyś się dowiedzieć o jego umiejętnościach?
            </span>
          </div>
  
          <!-- User messages with SVG icon on the right -->
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
  
        <!-- Input container for message entry -->
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
    </div>
  </template>
  
  <script>
  import '../../assets/css/Chat.css';
  
  export default {
    name: 'ChatApp',
    data() {
      return {
        messages: [],
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
  