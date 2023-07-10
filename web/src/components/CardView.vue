<template>
    <v-card :key="props.cardValue.code" class="mb-4 pb-3 flex-grow-1" @click="openCard" :color="getColor()">
        <v-tooltip activator="parent" location="right">
            <v-img :src="getImgURL" height="300" lazy-src="@/assets/marvel-player-back.png" transition="slide-x-transition" @error="onImgFail()"/>
        </v-tooltip>
        <v-card-title>{{ props.cardValue.card.name }}</v-card-title>
        <v-card-item class="pt-0 mt-n2">{{ props.cardValue.card.subname }}</v-card-item>
        <v-card-item class="mt-n3 text-h5">
            <b>{{ props.cardValue.value }}</b>
        </v-card-item>
        <v-chip>Popularity Mod: {{ props.cardValue.popularityMod.toFixed(3) }}</v-chip>
        <v-chip v-if="props.cardValue.eligibleHeroCount != 0">{{ getLockingTraits() }} Mod: {{ props.cardValue.traitMod.toFixed(3) }}</v-chip> 
    </v-card>
</template>

<script setup>
    import { ref, computed } from 'vue'

    const imgFailed = ref(false);
    const props = defineProps(['cardValue']);

    function openCard() {
        window.open("https://marvelcdb.com/card/" + props.cardValue.code);
    }

    function onImgFail() {
        imgFailed.value = true;
    }

    const getImgURL = computed(() => {
        return imgFailed.value ? "src/assets/marvel-player-back-not-found.png" : "https://marvelcdb.com/" + props.cardValue.card.imageSource;
    });

    function getColor() {
        if (props.cardValue.card.aspect == "basic") {
            return "grey-lighten-3"
        }
        if (props.cardValue.card.aspect == "aggression") {
            return "red-lighten-4"
        }
        if (props.cardValue.card.aspect == "protection") {
            return "green-lighten-4"
        }
        if (props.cardValue.card.aspect == "justice") {
            return "yellow-lighten-4"
        }
        if (props.cardValue.card.aspect == "leadership") {
            return "blue-lighten-4"
        }
    }

    function getLockingTraits() {
        if (props.cardValue.card.lockingTraits.length > 0) {
            let lts = props.cardValue.card.lockingTraits
            let s = lts[0]
            s = s.charAt(0).toUpperCase() + s.slice(1).toLowerCase()
            for (let i = 1; i < lts.length; i++) {
                let ns = lts[i]
                s += " or " + ns.charAt(0).toUpperCase() + ns.slice(1).toLowerCase()
            }
            return s
        }
    }
</script>
