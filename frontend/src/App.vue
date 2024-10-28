<template>
  <div id="app">
    <!-- Navigation Section -->
    <NavigationBar />

        <!-- Start Section -->
        <section class="section" id="start">
      <Start />
    </section>

    <!-- About Me Section -->
    <section class="section" id="about">
      <AboutMe />
    </section>

    <!-- Contact Form Section -->
    <section class="section" id="chat">
      <ChatApp />
    </section>

    <!-- Footer Section -->
    <section class="section" id="experience">
      <MyExperience />
    </section>

    <!-- Experience Section -->
    <section class="section" id="tech-stack">
      <TechStack />
    </section>

    

    <!-- Contact Form Section -->
    <section class="section" id="contact">
      <ContactForm />
    </section>

    <!-- Footer Section -->
    <section class="section" id="footer">
      <AppFooter />
    </section>
  </div>
</template>

<script>
import AboutMe from './components/Home/AboutMe.vue';
import TechStack from './components/Home/TechStack.vue';
import ContactForm from './components/Home/ContactForm.vue';
import NavigationBar from './components/Home/NavigationBar.vue';
import AppFooter from './components/Home/Footer.vue';
import MyExperience from './components/Home/Experience.vue';
import ChatApp from './components/Home/Chat.vue';
import Start from './components/Home/Start.vue';
import './assets/css/App.css';

export default {
  name: 'App',
  components: {
    AboutMe,
    TechStack,
    ContactForm,
    NavigationBar,
    AppFooter,
    MyExperience,
    ChatApp,
    Start,
  },
  mounted() {
    window.scrollTo(0, 0); // Resetuje scroll na górę strony po załadowaniu
    document.body.style.overflowX = 'hidden'; // Blokuje przewijanie w lewo i w prawo
  },
};
</script>

<style scoped>
/* Możesz dodać dodatkowe style tutaj */
</style>
