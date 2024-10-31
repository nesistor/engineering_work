<template>
  <div class="profile-background">
    <section id="about-me" class="about-me">
      <h1>My Bite Adventure</h1>
      <p>
        My life is a mix of culture, knowledge, and innovative ideas. I enjoy discovering new tools. Feel free to see me as a DevOps, Backend Developer, or Flutter Developer.
        I am very attentive to detail and consider the many factors necessary to prevent future issues. I strive for efficiency in all I do.
      </p>
    </section>
    
    <!-- New section with an image above the text -->
    <section id="new-section" class="new-section">
      <img src="@/assets/images/wolf.svg" alt="Image Description" class="section-image" />
      <h2>New Adventures Await</h2>
      <p>
        This is where I dive into my next projects and explore new territories in tech. I am always on the lookout for groundbreaking solutions and strive to bring innovative projects to life.
      </p>
    </section>
  </div>
</template>

<script>
import '../../assets/css/AboutMe.css';

export default {
  name: 'AboutMe',
  methods: {
    triggerAnimation() {
      const aboutSection = document.querySelector('#about-me');
      const newSection = document.querySelector('#new-section');
      
      // Resetowanie animacji
      aboutSection.classList.remove('slide-in');
      newSection.classList.remove('slide-in');
      
      setTimeout(() => {
        aboutSection.classList.add('slide-in');
        newSection.classList.add('slide-in');
      }, 50);
    }
  },
  mounted() {
    this.triggerAnimation();

    const aboutLink = document.querySelector('a[href="#about-me"]');
    if (aboutLink) {
      aboutLink.addEventListener('click', this.triggerAnimation);
    }
  },
  beforeUnmount() {
    const aboutLink = document.querySelector('a[href="#about-me"]');
    if (aboutLink) {
      aboutLink.removeEventListener('click', this.triggerAnimation);
    }
  }
};
</script>