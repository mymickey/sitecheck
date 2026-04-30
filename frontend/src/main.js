import './style.css'
import '../bindings/github.com/wailsapp/wails/v3/internal/eventcreate.js'

import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'

createApp(App).use(createPinia()).mount('#app')
