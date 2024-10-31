<template>
  <section class="start-section">
    <div class="content-wrapper">
      <h1 class="main-heading">Welcome to My Page</h1>
      <p class="description">Discover more about my work and projects below.</p>
      <div class="button-container">
        <button class="cta-button" @click="scrollToAboutMe">Discover More</button>
        <button class="cta-button">Load Resume</button>
      </div>
    </div>
  </section>
</template>

<script>
export default {
  name: "StartApp",
  methods: {
  scrollToAboutMe() {
    const aboutMeSection = document.querySelector('#about-me');
    if (aboutMeSection) {
      // Oblicz pozycję sekcji o kilka pikseli wyżej
      const offset = 300; // Ustawienie przesunięcia w pikselach
      const elementPosition = aboutMeSection.getBoundingClientRect().top + window.scrollY;
      const offsetPosition = elementPosition - offset;

      window.scrollTo({
        top: offsetPosition,
        behavior: 'smooth' // Płynne przewijanie
      });
    }
  },
},

};
</script>


  <style src="../../assets/css/Start.css"></style>
  