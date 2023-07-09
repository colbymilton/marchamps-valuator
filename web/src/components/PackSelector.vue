<template>
    <v-card class="ma-8 py-5">
        <div @click="expand = !expand">
            <div>
                <v-label class="text-h6">Owned Packs</v-label>
            </div>
            <div>
                <v-label v-show="!expand">Click to expand!</v-label>
            </div>
        </div>
        <v-expand-transition>
            <div v-show="expand">
                <v-row dense class="pa-5" justify="center">
                    <v-col v-for="pack in store.packs" cols="12" sm="4" lg="2">
                        <div :key="pack.code" class="pa-0 mx-3 my-n3">
                            <v-checkbox :label="pack.name" v-model="pack.owned" :disabled="pack.locked" class="black"/>
                        </div>
                    </v-col>
                </v-row>
                <v-btn ripple color="secondary" :loading="loading" @click="getValues">Get Values</v-btn>
            </div>
        </v-expand-transition>

        <!-- buttons along the bottom to clear selections and also to select all -->
    </v-card>
</template>

<script setup>
    import { useAppStore } from '@/store/app';
    import { reactive, ref } from 'vue';

    const expand = ref(true);
    const store = reactive(useAppStore());
    const loading = ref(false);

    async function getValues() {
        loading.value = true;
        let s = store.packsString();
        const result = await fetch('http://localhost:9999/pack_values?owned=' + s);
        const data = await result.json();
        store.packValues = data;
        store.selectedPack = store.packValues[0]
        loading.value = false;
        expand.value = false;
    };
</script>
