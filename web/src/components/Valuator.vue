<template>
  <v-responsive class="align-center text-center">
      <PackSelector/>
      <PackViewer/>
  </v-responsive>
</template>
  
<script setup>
  import PackSelector from '@/components/PackSelector.vue'
  import PackViewer from '@/components/PackViewer.vue'
</script>
  