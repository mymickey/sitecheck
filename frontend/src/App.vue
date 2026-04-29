<script setup>
import { computed, onMounted, shallowRef } from 'vue'
import { Activity, Info } from 'lucide-vue-next'
import AboutPanel from './components/AboutPanel.vue'
import ConnectivitySettings from './components/ConnectivitySettings.vue'
import TrayMenu from './components/TrayMenu.vue'
import { useSiteCheckStore } from './stores/sitecheck'

const store = useSiteCheckStore()
const isMenuMode = window.location.search.includes('mode=menu')
const activeTab = shallowRef('connectivity')

const tabs = [
  { id: 'connectivity', label: 'Connectivity', icon: Activity },
  { id: 'about', label: 'About', icon: Info },
]

const activeComponent = computed(() => (
  activeTab.value === 'about' ? AboutPanel : ConnectivitySettings
))

onMounted(() => {
  store.loadSettings()
})
</script>

<template>
  <div :class="{ 'is-menu-mode': isMenuMode }">
    <TrayMenu v-if="isMenuMode" />
    <div v-else class="app-container">
      <header class="window-header"></header>
      <main class="app-shell default">
    <aside class="sidebar">
      <div class="brand">
        <img class="brand__mark" src="/logo.svg" alt="" aria-hidden="true">
        <div class="brand__copy">
          <strong>SiteCheck</strong>
          <span>macOS menubar</span>
        </div>
      </div>

      <nav class="tabs" data-orientation="vertical" aria-label="Settings sections">
        <div class="tabs__list" data-orientation="vertical">
          <button
            v-for="tab in tabs"
            :key="tab.id"
            class="tabs__tab nav-tab"
            type="button"
            :data-selected="activeTab === tab.id"
            @click="activeTab = tab.id"
          >
            <component :is="tab.icon" class="nav-tab__icon" aria-hidden="true" />
            <span>{{ tab.label }}</span>
          </button>
        </div>
      </nav>
    </aside>

    <section class="content" aria-live="polite">
      <component :is="activeComponent" />
    </section>
  </main>
  </div>
  </div>
</template>
