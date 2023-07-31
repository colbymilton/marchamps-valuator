<template>
    <v-btn icon v-bind="props">
        <v-icon>mdi-information</v-icon>
        <v-dialog v-model="dialog" activator="parent" width="auto" scrollable max-width="1100">
            <v-card class="align-center pa-4" color="secondary" border>
                <v-card-title>What is the Marvel Champions Pack Valuator?</v-card-title>
                <v-card-text>
                    This tool is meant to help Marvel Champions players 
                    decide which pack they should buy next when attempting to build up their collection. 
                    It's worth noting that I believe players should primarily focus on acquiring the 
                    heroes (and villains) that interest them - but if you've already bought your favorite heroes, 
                    this tool might help you decide what to get next!

                    <br/><br/>

                    Please note that this tool is focused on improving your collection for deck building. 
                    That means packs are not evaluated based on the hero that they include; 
                    they are instead evaluated based on the <i>other</i> cards in the pack, as those can typically
                    be used with any other hero. 

                    <br/><br/>

                    <i>
                    For example: Hulk is widely considered to be a weak hero but his pre-constructed deck 
                    includes some fantastic aggression cards. Due to these great included cards, 
                    Hulk's pack may end up being valued quite highly by this tool despite the fact Hulk himself is typically not.
                    </i>

                    <br/><br/>

                    For more information about the tool including details on how the cards are evaluated, please check out the
                    Github repository that is linked below!
                </v-card-text>
            </v-card>
        </v-dialog>
    </v-btn>
</template>

<script setup>
    import { ref } from 'vue'

    const dialog = ref(false);
</script>
