<template>
    <v-footer color="secondary" border>
      <v-row justify="center" class="my-2 mx-5 text-center text-body-3" dense>
        <v-col cols="12">
          This web app is made possible by the API that is provided by the wonderful <a href="https://marvelcdb.com/">Marvelcdb.com</a>.<br/>
          I thank the team behind marvelcdb and also the many users who use the service and build fun decks!
        </v-col>
        <v-divider/>
        <v-col cols="12">
          This web app does not claim ownership of Marvel Champions or its card artwork/information.<br/>
          All game-related content belongs to Fantasy Flight Games. This is purely for informational purposes.<br/>
          Please support the game at <a href="https://www.fantasyflightgames.com/en/products/marvel-champions-the-card-game/">Fantasy Flight Games</a>.
        </v-col>
        <v-divider/>
        <v-col cols="12">
          This project is open source and available on my <a href="https://github.com/colbymilton/marchamps-valuator">Github</a>.<br/>
          Creating Github issues for feature requests and bug reports is welcome!
        </v-col>
        <v-divider/>
        <v-col cols="12">
          This tool was made by a solo dev for fun and the love of Marvel Champions.<br/>
          If you find it useful, consider donating to fund my card game addiction: <a href="https://bmc.link/miltoncolbydev">Donate</a>.
        </v-col>
        <v-col cols="12">
          <v-row justify="center" class="ma-3">
            <v-tooltip text="Marvelcdb" location="top">
              <template v-slot:activator="{ props }">
                <v-btn v-bind="props" icon class="mx-2" href="https://marvelcdb.com/">
                  <v-icon>mdi-database-outline</v-icon>
                </v-btn>              
              </template>
            </v-tooltip>
            <v-tooltip text="Fantasy Flight Games" location="top">
              <template v-slot:activator="{ props }">
                <v-btn v-bind="props" icon class="mx-2" href="https://www.fantasyflightgames.com/en/products/marvel-champions-the-card-game/">
                  <v-icon>mdi-alpha-f-box-outline</v-icon>
                </v-btn>              
              </template>
            </v-tooltip>
            <v-tooltip text="Github" location="top">
              <template v-slot:activator="{ props }">
                <v-btn v-bind="props" icon class="mx-2" href="https://github.com/colbymilton/marchamps-valuator">
                  <v-icon>mdi-github</v-icon>
                </v-btn>              
              </template>
            </v-tooltip>
            <v-tooltip text="Donation" location="top">
              <template v-slot:activator="{ props }">
                <v-btn v-bind="props" icon class="mx-2" href="https://bmc.link/miltoncolbydev">
                  <v-icon>mdi-hand-coin</v-icon>
                </v-btn>              
              </template>
            </v-tooltip>
          </v-row>
        </v-col>
      </v-row>
    </v-footer>
</template>