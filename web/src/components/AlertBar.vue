<template>
    <v-alert v-model="hasError" icon="false" closable type="error" @click:close="closeError()">{{ store.error }}</v-alert>
</template>

<script setup>
    import { computed } from 'vue'
    import { useAppStore } from '@/store/app';

    const store = useAppStore();

    const hasError = computed(() => {
        return store.error != "";
    });

    function closeError() {
        console.log("closed error:", store.error);
        store.error = "";
    }
</script>