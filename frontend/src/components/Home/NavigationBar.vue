<template>
  <nav>
    <div class="navbar-logo">
      <img src="@/assets/images/logo_km.svg" alt="Logo" class="logo-icon" />
      <span class="name-text">Karol Malicki</span> <!-- Zmiana na zwykły tekst -->
    </div>
    <ul>
      <li><a href="#about">About Me</a></li>
      <li><a href="#tech-stack">Tech Stack</a></li>
      <li><a href="https://github.com/nesistor" target="_blank" rel="noopener noreferrer">Github</a></li>
      <li><a href="#contact" class="contact-link">Contact</a></li>
    </ul>
  </nav>
</template>

<script>
import '../../assets/css/NavigationBar.css';

export default {
  name: 'NavigationBar',
};
</script>
