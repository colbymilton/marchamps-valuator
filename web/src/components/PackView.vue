<template>
    <v-card :key="props.packValue.code" :title="props.packValue.pack.name" :color="props.packValue.code == store.selectedPack.code? 'amber-lighten-3' : 'amber-lighten-4'" class="mb-4 pb-4"  @click="selectPack">
        <v-card-item class="mt-n3 text-h5">
            <b>{{ props.packValue.valueSum }}</b>
        </v-card-item>
        <v-chip>New cards: {{ countNew() }}</v-chip>
    </v-card>
</template>

<script setup>
    import { useAppStore } from '@/store/app';
    import { reactive } from 'vue';

    const props = defineProps(['packValue'])
    const store = reactive(useAppStore());

    async function selectPack() {
        store.selectedPack = props.packValue;
        window.scrollTo(0, 0);
    }

    function countNew() {
        let count = 0;
        for (let i = 0; i < props.packValue.cardValues.length; i++) {
            let cv = props.packValue.cardValues[i]
            if (cv.newMod == 1) {
                count++
            }
        }
        return count;
    }

</script>
