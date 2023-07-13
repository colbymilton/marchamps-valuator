<template>
    <v-card class="ma-8 py-5">
        <div @click="expandOptions = !expandOptions">
            <div>
                <v-label class="text-h5">Options</v-label>
            </div>
            <div>
                <v-label v-show="!expandOptions">Click to expand!</v-label>
            </div>
        </div>
        <v-expand-transition>
            <div v-show="expandOptions">
                <div>
                    <v-label class="text-h6">Owned Packs</v-label>
                </div>
                <v-progress-circular v-if="loadingPacks" indeterminate class="my-3"/>
                <v-row v-if="!loadingPacks" dense class="pa-5" justify="center">
                    <v-col v-for="pack in store.packs" cols="12" sm="4" lg="2">
                        <div :key="pack.code" class="pa-0 mx-3 my-n6">
                            <v-checkbox :label="pack.name" v-model="pack.owned" :disabled="pack.locked" class="black"/>
                        </div>
                    </v-col>
                </v-row>

                <div>
                    <v-label class="text-h6">Advanced</v-label>
                </div>
                <div class="mt-2">
                    <v-label>Card Aspect Weights</v-label>
                    <v-row class="mx-6 mt-1">
                        <v-col cols="12" md="6" v-for="weight in weights">
                            <v-row>
                                <v-col cols="2"><v-label class="mt-1">{{ weight.aspect }}</v-label></v-col>
                                <v-col cols="10"><v-slider min="0" max="1" show-ticks step="0.05" thumb-label v-model="weight.weight" :color="getColor(weight.aspect)"/></v-col>
                            </v-row>
                        </v-col>
                    </v-row>
                </div>
                <!--
                <div class="mt-2">
                    <v-label>Pack Inclusions</v-label>
                    <v-row dense class="mx-6" justify="center">
                        <v-col cols="2"><v-checkbox label="Include Campaign Expansions" v-model="includeCampaigns"/></v-col>
                        <v-col cols="2"><v-checkbox label="Include Heroes" v-model="includeHeroes"/></v-col>
                    </v-row>
                </div>
                -->
                <v-btn ripple color="secondary" :loading="loadingValues" @click="getValues">Get Values</v-btn>
            </div>
        </v-expand-transition>

        <!-- buttons along the bottom to clear selections and also to select all -->
    </v-card>
</template>

<script setup>
    import { useAppStore } from '@/store/app';
    import { ref, reactive } from 'vue';

    const store = useAppStore();

    const loadingValues = ref(false);
    const loadingPacks = ref(true);
    const expandOptions = ref(true);

    const weights = reactive([
        {
            aspect: "Aggression",
            code: "aw",
            weight: 1,
        },
        {
            aspect: "Protection",
            code: "pw",
            weight: 1,
        },
        {
            aspect: "Justice",
            code: "jw",
            weight: 1,
        },
        {
            aspect: "Leadership",
            code: "lw",
            weight: 1,
        },
    ]);

    const includeCampaigns = ref(true);
    const includeHeroes = ref(true);

    // load packs
    fetch('http://localhost:9999/packs')
        .then(async response => {
            const data = await response.json();
            if (response.ok) {
                store.packs = data;
                for (let i = 0; i < store.packs.length; i++) {
                    let pack = store.packs[i];
                    if (pack.code == "core") {
                        pack.owned = true;
                        pack.locked = true;
                    }
                }
            } else {
                store.error = data.error;
            }
        })
        .catch(error => {
            store.error = error;
        })
        .finally(() => {
            loadingPacks.value = false;
        });

    async function getValues() {
        loadingValues.value = true;
        let s = store.packsString();

        // get weights
        let weightStr = ""
        for (let i = 0; i < weights.length; i++) {
            let weight = weights[i];
            weightStr += "&" + weight.code + "=" + weight.weight;
        }

        fetch('http://localhost:9999/pack_values?owned=' + s + weightStr)
            .then(async response => {
                const data = await response.json();
                if (response.ok) {
                    store.packValues = data;
                    store.selectedPack = store.packValues[0];
                    expandOptions.value = false;
                } else {
                    store.error = data.error;
                }
            })
            .catch(error => {
                store.error = error;
            })
            .finally(() => {
                loadingValues.value = false;
            });        
    };

    function getColor(aspect) {
        if (aspect == "Basic") {
            return "grey-lighten-3"
        }
        if (aspect == "Aggression") {
            return "red-lighten-4"
        }
        if (aspect == "Protection") {
            return "green-lighten-4"
        }
        if (aspect == "Justice") {
            return "yellow-lighten-4"
        }
        if (aspect == "Leadership") {
            return "blue-lighten-4"
        }
    }
</script>
