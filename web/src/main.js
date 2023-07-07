/**
 * main.js
 *
 * Bootstraps Vuetify and other plugins then mounts the App`
 */

import App from './App.vue'
import { createApp } from 'vue'
import { registerPlugins } from '@/plugins'
import { createPinia } from 'pinia';

const app = createApp(App).use(createPinia())

registerPlugins(app)

app.mount('#app')
