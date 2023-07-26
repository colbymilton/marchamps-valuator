<template>
  <v-app-bar color="primary">
    <v-app-bar-title>
      <v-icon icon="mdi-cards" />
      Marvel Champions Pack Valuator!
    </v-app-bar-title>
    <about-popup/>
    <theme-selector/>
  </v-app-bar>
</template>

<script setup>
  import AboutPopup from '@/components/AboutPopup.vue';
  import ThemeSelector from '@/components/ThemeSelector.vue';
</script>
