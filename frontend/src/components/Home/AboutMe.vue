<template>
  <section id="about-me" class="about-me">
    <h1>My bite adventure</h1>
    <p>My life is a mix of culture, knowledge, and innovative ideas. I enjoy discovering new tools. Feel free to see me as a DevOps, Backend Developer, or Flutter Developer.
      I am very attentive to detail and consider the many factors necessary to prevent future issues. I strive for efficiency in all I do.</p>
  </section>
</template>

<script>
import '../../assets/css/AboutMe.css';

export default {
  name: 'AboutMe',
  methods: {
    triggerAnimation() {
      const aboutSection = document.querySelector('#about-me');
      aboutSection.classList.remove('slide-in'); // Usuwamy klasę, by zresetować animację
      // Ustawiamy ponownie klasę po krótkim czasie
      setTimeout(() => {
        aboutSection.classList.add('slide-in');
      }, 50); // Krótkie opóźnienie dla płynności efektu
    }
  },
  mounted() {
    // Animacja przy załadowaniu komponentu
    this.triggerAnimation();

    // Dodajemy nasłuchiwanie kliknięcia na link About Me w navbarze
    const aboutLink = document.querySelector('a[href="#about-me"]');
    if (aboutLink) {
      aboutLink.addEventListener('click', this.triggerAnimation);
    }
  },
  beforeUnmount() {
    // Usuwamy nasłuchiwanie zdarzenia przy usuwaniu komponentu
    const aboutLink = document.querySelector('a[href="#about-me"]');
    if (aboutLink) {
      aboutLink.removeEventListener('click', this.triggerAnimation);
    }
  }
};
</script>
