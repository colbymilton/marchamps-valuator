<template>
    <v-card :key="props.packValue.code" class="mb-4 pb-4"
    :color="props.packValue.code == store.selectedPack.code? 'accent' : ''" @click="selectPack">
        <v-card-title class="mx-n3 text-wrap">{{ props.packValue.pack.name }}</v-card-title>
        <v-card-item class="mt-n3 text-h5">
            <b>{{ props.packValue.valueSum }}</b>
        </v-card-item>
        <v-chip>New cards: {{ countNew }}</v-chip>
    </v-card>
</template>

<script setup>
    import { useAppStore } from '@/store/app';
    import { computed } from 'vue';

    const props = defineProps(['packValue']);
    const store = useAppStore();

    async function selectPack() {
        store.selectedPack = props.packValue;
        window.scrollTo(0, 0);
    }

    const countNew = computed(() => {
        let count = 0;
        for (let i = 0; i < props.packValue.cardValues.length; i++) {
            let cv = props.packValue.cardValues[i];
            if (cv.newMod == 1) {
                count++;
            }
        }
        return count;
    });

</script>
