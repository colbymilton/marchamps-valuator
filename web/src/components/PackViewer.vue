<template>
    <v-row class="mx-8">
        <v-col cols="4">
            <PackView v-for="packValue in store.packValues" :pack-value="packValue"/>
        </v-col>
        <v-col cols="8">
            <v-row no-gutters>
                <v-col v-for="cardValue in store.selectedPack.cardValues" cols="12" sm="6" md="4" lg="3" class="px-2 d-flex flex-column">
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
