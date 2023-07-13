<template>
  <v-responsive class="align-center text-center">
      <AlertBar style="position:fixed;z-index:1000"/>
      <PackSelector/>
      <PackViewer/>
  </v-responsive>
</template>
  
<script setup>
  import PackSelector from '@/components/PackSelector.vue'
  import PackViewer from '@/components/PackViewer.vue'
  import AlertBar from '@/components/AlertBar.vue'
</script>
  