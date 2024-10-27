<template>
  <div id="app">
    <!-- Navigation Section -->
    <NavigationBar />

    <!-- About Me Section -->
    <section class="section" id="about">
      <AboutMe />
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
import AboutMe from './components/AboutMe.vue';
import TechStack from './components/TechStack.vue';
import ContactForm from './components/ContactForm.vue';
import NavigationBar from './components/NavigationBar.vue';
import AppFooter from './components/Footer.vue';
import MyExperience from './components/Experience.vue';
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
