<template>
  <v-app>
    <default-bar />

    <default-view />

    <v-spacer/>
    <v-spacer/>
    <v-spacer/>
    <default-footer/>
  </v-app>
</template>

<script setup>
  import DefaultBar from './AppBar.vue'
  import DefaultView from './View.vue'
  import DefaultFooter from './Footer.vue'
</script>
