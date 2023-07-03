<template>
    <v-card title="Owned Packs" color="blue-lighten-4" class="ma-8">
        <v-row dense class="pa-5" justify="center">
            <div v-for="pack in packs" :key="pack.code" class="pa-0 mx-3 my-n3">
                <v-checkbox :label="pack.name"/>
            </div>
        </v-row>

        <!-- buttons along the bottom to clear selections and also to select all -->
    </v-card>
</template>

<script setup>
import { reactive } from 'vue';

    let js = `[
    {
        "code": "core",
        "name": "Core Set",
        "id": 1,
        "available": "2019-11-01"
    },
    {
        "code": "gob",
        "name": "The Green Goblin",
        "id": 2,
        "available": "2019-12-20"
    },
    {
        "code": "cap",
        "name": "Captain America",
        "id": 3,
        "available": "2019-12-20"
    },
    {
        "code": "msm",
        "name": "Ms. Marvel",
        "id": 4,
        "available": "2019-12-20"
    },
    {
        "code": "twc",
        "name": "The Wrecking Crew",
        "id": 6,
        "available": "2020-02-07"
    },
    {
        "code": "thor",
        "name": "Thor",
        "id": 5,
        "available": "2020-03-06"
    },
    {
        "code": "bkw",
        "name": "Black Widow",
        "id": 7,
        "available": "2020-05-06"
    },
    {
        "code": "drs",
        "name": "Doctor Strange",
        "id": 8,
        "available": "2020-06-06"
    },
    {
        "code": "ron",
        "name": "Ronan Modular Set",
        "id": 16,
        "available": "2020-08-02"
    },
    {
        "code": "hlk",
        "name": "Hulk",
        "id": 9,
        "available": "2020-08-07"
    },
    {
        "code": "trors",
        "name": "The Rise of Red Skull",
        "id": 10,
        "available": "2020-09-04"
    },
    {
        "code": "toafk",
        "name": "The Once and Future Kang",
        "id": 11,
        "available": "2020-10-02"
    },
    {
        "code": "ant",
        "name": "Ant-man",
        "id": 12,
        "available": "2020-11-06"
    },
    {
        "code": "wsp",
        "name": "Wasp",
        "id": 13,
        "available": "2021-01-22"
    },
    {
        "code": "qsv",
        "name": "Quicksilver",
        "id": 14,
        "available": "2021-02-05"
    },
    {
        "code": "scw",
        "name": "Scarlet Witch",
        "id": 15,
        "available": "2021-03-05"
    },
    {
        "code": "gmw",
        "name": "Galaxy's Most Wanted",
        "id": 17,
        "available": "2021-04-02"
    },
    {
        "code": "stld",
        "name": "Star-Lord",
        "id": 18,
        "available": "2021-05-14"
    },
    {
        "code": "gam",
        "name": "Gamora",
        "id": 19,
        "available": "2021-05-14"
    },
    {
        "code": "drax",
        "name": "Drax",
        "id": 20,
        "available": "2021-06-18"
    },
    {
        "code": "vnm",
        "name": "Venom",
        "id": 21,
        "available": "2021-07-16"
    },
    {
        "code": "mts",
        "name": "The Mad Titan's Shadow",
        "id": 22,
        "available": "2021-09-01"
    },
    {
        "code": "nebu",
        "name": "Nebula",
        "id": 23,
        "available": "2021-09-17"
    },
    {
        "code": "warm",
        "name": "War Machine",
        "id": 24,
        "available": "2021-11-05"
    },
    {
        "code": "hood",
        "name": "The Hood",
        "id": 25,
        "available": "2021-11-26"
    },
    {
        "code": "valk",
        "name": "Valkyrie",
        "id": 26,
        "available": "2021-11-26"
    },
    {
        "code": "vision",
        "name": "Vision",
        "id": 27,
        "available": "2022-01-14"
    },
    {
        "code": "sm",
        "name": "Sinister Motives",
        "id": 28,
        "available": "2022-04-08"
    },
    {
        "code": "nova",
        "name": "Nova",
        "id": 29,
        "available": "2022-05-20"
    },
    {
        "code": "ironheart",
        "name": "Ironheart",
        "id": 30,
        "available": "2022-05-20"
    },
    {
        "code": "spiderham",
        "name": "Spider-Ham",
        "id": 31,
        "available": "2022-07-15"
    },
    {
        "code": "spdr",
        "name": "SP//dr",
        "id": 32,
        "available": "2022-07-15"
    },
    {
        "code": "mut_gen",
        "name": "Mutant Genesis",
        "id": 33,
        "available": "2022-09-30"
    },
    {
        "code": "cyclops",
        "name": "Cyclops",
        "id": 34,
        "available": "2022-09-30"
    },
    {
        "code": "phoenix",
        "name": "Phoenix",
        "id": 35,
        "available": "2022-09-30"
    },
    {
        "code": "wolv",
        "name": "Wolverine",
        "id": 37,
        "available": "2022-11-11"
    },
    {
        "code": "storm",
        "name": "Storm",
        "id": 38,
        "available": "2022-11-11"
    },
    {
        "code": "mojo",
        "name": "Mojo Mania",
        "id": 36,
        "available": "2022-11-11"
    },
    {
        "code": "gambit",
        "name": "Gambit",
        "id": 39,
        "available": "2023-02-24"
    },
    {
        "code": "rogue",
        "name": "Rogue",
        "id": 40,
        "available": "2023-02-24"
    }
]`

    const packs = reactive(JSON.parse(js))

</script>
