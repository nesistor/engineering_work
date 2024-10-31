<template>
  <div class="experience-background">
    <section class="experience">
      <h1>My Bite Adventure</h1>
      <div
        class="experience-item"
        v-for="(experience, index) in experiences"
        :key="index"
        @mousemove="handleMouseMove"
        @mouseleave="handleMouseLeave"
        @mouseover="handleMouseOver"
        ref="experienceItem"
      >
        <div class="experience-header">
          <span class="company-name">{{ experience.company }}</span>
          <span class="employment-dates">{{ experience.dates }}</span>
        </div>
        <p class="job-description">{{ experience.description }}</p>
        <div class="glow" />
      </div>
    </section>
  </div>
</template>

<script>
export default {
  name: "MyExperience",
  data() {
    return {
      experiences: [
        {
          company: "Tech Solutions Inc.",
          dates: "Jan 2020 - Dec 2022",
          description: `- Led a team to implement CI/CD pipelines and automation.
            - Collaborated with cross-functional teams to improve efficiency.
            - Developed backend services and APIs in Node.js and GoLang.`,
        },
        {
          company: "Innovatech Corp.",
          dates: "Feb 2018 - Dec 2019",
          description: `- Designed and maintained cloud infrastructure solutions.
            - Optimized application performance and reliability.
            - Provided training for new hires on DevOps best practices.`,
        },
      ],
    };
  },
  methods: {
    handleMouseOver(event) {
      const glow = event.currentTarget.querySelector(".glow");
      glow.style.opacity = 1;
    },
    handleMouseLeave(event) {
      const card = event.currentTarget;
      const glow = card.querySelector(".glow");
      card.style.transform = `perspective(500px) scale(1) rotateX(0) rotateY(0)`;
      glow.style.opacity = 0;
    },
    handleMouseMove(event) {
      const card = event.currentTarget;
      const relX = (event.offsetX + 1) / card.offsetWidth;
      const relY = (event.offsetY + 1) / card.offsetHeight;
      const rotY = `rotateY(${(relX - 0.5) * 30}deg)`;
      const rotX = `rotateX(${(relY - 0.5) * -30}deg)`;
      card.style.transform = `perspective(500px) scale(1.05) ${rotY} ${rotX}`;
    },
  },
};
</script>

<style scoped src="../../assets/css/Experience.css"></style>
