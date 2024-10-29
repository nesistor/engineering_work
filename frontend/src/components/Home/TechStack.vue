<template>
  <section class="tech-stack">
    <h1>Software</h1>
    <ul>
      <!-- Render each technology item with its icon only -->
      <li v-for="(icon, index) in techIcons" :key="index">
        <img :src="icon" alt="Technology Icon" class="tech-icon" />
      </li>
    </ul>
  </section>
</template>

<script>
import '../../assets/css/TechStack.css';

export default {
  name: 'TechStack',
  props: {
    techIcons: {
      type: Array,
      default: () => [
        require('@/assets/images/soft_icons/golang.svg'),
        require('@/assets/images/soft_icons/vue.svg'),
        require('@/assets/images/soft_icons/python.svg'),
        require('@/assets/images/soft_icons/aws.svg'),
        require('@/assets/images/soft_icons/flutter.svg'),
        // Add more icon paths as needed
      ]
    }
  }
}
</script>
