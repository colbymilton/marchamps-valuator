<template>
    <v-row class="mx-8">
        <v-col>
            <PackView v-for="packValue in store.packValues" :pack-value="packValue"/>
        </v-col>
        <v-col>
            <v-row no-gutters>
                <v-col v-for="cardValue in store.selectedPack.cardValues" cols="6" class="px-2">
                    <CardView :card-value="cardValue"/>
                </v-col>
            </v-row>
        </v-col>
    </v-row>
</template>

<script setup>
    import PackView from '@/components/PackView.vue'
    import CardView from '@/components/CardView.vue'
    import { useAppStore } from '@/store/app';
    import { reactive } from 'vue';

    const store = reactive(useAppStore());
</script>
