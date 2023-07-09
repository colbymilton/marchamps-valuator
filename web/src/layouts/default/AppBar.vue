<template>
  <v-app-bar color="primary">
    <v-app-bar-title>
      <v-icon icon="mdi-cards" />
      Marvel Champions Pack Valuator!
    </v-app-bar-title>
    <v-spacer/>
    <theme-selector/>
  </v-app-bar>
</template>

<script setup>
  import ThemeSelector from '@/components/ThemeSelector.vue';
</script>
