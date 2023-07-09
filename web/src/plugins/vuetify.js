/**
 * plugins/vuetify.js
 *
 * Framework documentation: https://vuetifyjs.com`
 */

// Styles
import '@mdi/font/css/materialdesignicons.css'
import 'vuetify/styles'

// Composables
import { createVuetify } from 'vuetify'

// colors: {
//   background ,
//   surface ,
//   primary: ,
//   secondary: ,
//   accent: ,
//   error: ,
//   info: ,
//   success: ,
//   warning: ,
// }


// https://vuetifyjs.com/en/introduction/why-vuetify/#feature-guides
export default createVuetify({
  theme: {
    defaultTheme: 'captainAmerica',
    themes: {
      blackPanther: {
        dark: true,
        colors: {
            primary: '#493C90',
            secondary: '#843690',
            accent: '#FFE487',
          }
      },
      ironMan: {
        dark: false,
        colors: {
            background: '#F6F6F6',
            primary: '#77160E',
            secondary: '#E30819',
            accent: '#FEDD12',
          }
      },
      spiderMan: {
        dark: false,
        colors: {
            background: '#1873B8',
            primary: '#BF1724',
            secondary: '#E30819',
            accent: '#00497C',
          }
      },
      scarletWitch: {
        dark: true,
        colors: {
            primary: '#E50719',
            secondary: '#B91444',
            accent: '#D60B54',
          }
      },
      loki: {
        dark: true,
        colors: {
            primary: '#008D39',
            secondary: '#016634',
            accent: '#F3D301',
          }
      },
      captainAmerica: {
        dark: false,
        colors: {
            background: '#ECEDEA',
            primary: '#015198',
            secondary: '#0369B3',
            accent: '#E30819',
          }
      }
    },
  },
})
