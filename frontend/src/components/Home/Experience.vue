<template>
  <div class="experience-background">
    <section class="experience">
      <h1>My Bite Adventure</h1>

      <!-- Pojedyncza pozycja doświadczenia zawodowego -->
      <div class="experience-item">
        <div class="experience-header">
          <span class="company-name">Tech Solutions Inc.</span>
          <span class="employment-dates">Jan 2020 - Dec 2022</span>
        </div>
        <p class="job-description">
          - Led a team to implement CI/CD pipelines and automation.<br />
          - Collaborated with cross-functional teams to improve efficiency.<br />
          - Developed backend services and APIs in Node.js and GoLang.
        </p>
      </div>

      <!-- Dodaj więcej doświadczeń zawodowych w tym samym formacie -->
      <div class="experience-item">
        <div class="experience-header">
          <span class="company-name">Innovatech Corp.</span>
          <span class="employment-dates">Feb 2018 - Dec 2019</span>
        </div>
        <p class="job-description">
          - Designed and maintained cloud infrastructure solutions.<br />
          - Optimized application performance and reliability.<br />
          - Provided training for new hires on DevOps best practices.
        </p>
      </div>
    </section>
  </div>
</template>

<script>
import '../../assets/css/Experience.css';
export default {
  name: 'MyExperience',
};
</script>
