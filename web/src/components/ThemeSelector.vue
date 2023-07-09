<template>
    <v-menu>
        <template v-slot:activator="{ props }">
            <v-btn icon v-bind="props">
                <v-icon>mdi-theme-light-dark</v-icon>
            </v-btn>
        </template>
        <v-card>
            <v-row>
                <v-col cols="6">
                    <v-label class="ml-4 text-h6">Light Themes</v-label>
                    <v-list-item v-for="theme in light_themes" :key="theme" @click="changeTheme(theme)">
                        <v-label>{{ theme.title }}</v-label> 
                    </v-list-item>
                </v-col>
                <v-col cols="6">
                    <v-label class="ml-4 text-h6">Dark Themes</v-label>
                    <v-list-item v-for="theme in dark_themes" :key="theme" @click="changeTheme(theme)">
                        <v-label>{{ theme.title }}</v-label> 
                    </v-list-item>
                </v-col>
            </v-row>
        </v-card>

    </v-menu>

</template>

<script setup>
    import { useTheme } from 'vuetify'
    
    const vuetifyTheme = useTheme()

    let dark_themes = [
        {
            title: "Scarlet Witch",
            code: "scarletWitch",
        },
        {
            title: "Black Panther",
            code: "blackPanther",
        },
        {
            title: "Loki",
            code: "loki",
        }
    ]

    let light_themes = [
        {
            title: "Captain America",
            code: "captainAmerica",
        },
        {
            title: "Iron Man",
            code: "ironMan",
        },
        {
            title: "Spider-Man",
            code: "spiderMan",
        },
    ]

    function changeTheme(theme) {
        console.log("pressed", theme)
        vuetifyTheme.global.name.value = theme.code;
    }
</script>
