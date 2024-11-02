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
          <!-- Wyświetlanie unikalnej ikony dla każdej pracy -->
          <img :src="experience.icon" alt="Company Icon" class="company-icon" />
          <span class="company-name">{{ experience.company }}</span>
          <span class="employment-dates">{{ experience.dates }}</span>
        </div>
        <!-- Opis pracy -->
        <p class="job-description" v-html="experience.description"></p>
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
          company: "Vocale Sp.z.o.o",
          dates: "Jan 2023 - Jun 2023",
          icon: require('@/assets/images/company/1.svg'), // Ścieżka do ikony Vocale
          description: `At this company, I developed a complete application from start to finish.<br><br>
          - I created responsive widgets that matched Figma designs precisely.<br>
          - I built a chat page, presenting state management solutions like Bloc and Provider to the team.<br>
          - My role involved extensive application testing to ensure quality.<br>
          - This experience strengthened my ability to deliver high-quality applications.`,
        },
        {
          company: "Marchesini Group S.p.A",
          dates: "Jun 2024 - Jun 2024",
          icon: require('@/assets/images/company/2.svg'), // Ścieżka do ikony Marchesini
          description: `At this company, I had the opportunity to collaborate with my team to create a comprehensive pitch deck and a demo product.<br><br>
          - Developed innovative concepts for automated pharmaceutical packaging machinery.<br>
          - Demonstrated the use of Document AI combined with NLP techniques to process prescription data and extract relevant information for storage in a database.<br>
          - Designed a program in C++ for Arduino that controls the display of specific LEDs at designated times.<br>
          - Applied Industrial Innovation techniques for Idea Validation Conducted Market Research and Competitor Analysis.<br>
          - Created and presented a Pitch Deck for the final evaluation.`,
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
