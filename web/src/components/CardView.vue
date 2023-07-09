<template>
    <v-card :key="props.cardValue.code" :title="props.cardValue.card.name" class="pb-4 mb-4 flex-grow-1" @click="openCard" :color="getColor()">
        <v-tooltip activator="parent" location="top">{{ props.cardValue.card.name }}</v-tooltip>
        <v-card-item class="pt-0 mt-n3">{{ props.cardValue.card.subname }}</v-card-item>
        <v-spacer/>
        <v-card-item class="mt-n3 text-h5">
            <b>{{ props.cardValue.value }}</b>
        </v-card-item>
        <v-chip>Popularity Mod: {{ props.cardValue.popularityMod.toFixed(3) }}</v-chip>
        <v-chip v-if="props.cardValue.eligibleHeroCount != 0">{{ getLockingTraits() }} Mod: {{ props.cardValue.traitMod.toFixed(3) }}</v-chip> 
    </v-card>
</template>

<script setup>
    const props = defineProps(['cardValue'])

    function openCard() {
        window.open("https://marvelcdb.com/card/" + props.cardValue.code)
    }

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
