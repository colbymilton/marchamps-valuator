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
                    <v-list-item v-for="theme in lightThemes" :key="theme" @click="changeTheme(theme.code)">
                        <v-label>{{ theme.title }}</v-label> 
                    </v-list-item>
                </v-col>
                <v-col cols="6">
                    <v-label class="ml-4 text-h6">Dark Themes</v-label>
                    <v-list-item v-for="theme in darkThemes" :key="theme" @click="changeTheme(theme.code)">
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

    let darkThemes = [
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

    let lightThemes = [
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

    // on setup load theme from localstorage
    let storedTheme = localStorage.getItem("marchamps-valuator-theme")
    if (storedTheme != null) {
        changeTheme(storedTheme)
    }

    function changeTheme(code) {
        vuetifyTheme.global.name.value = code;
        localStorage.setItem("marchamps-valuator-theme", code)
    }
</script>
