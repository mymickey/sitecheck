<script setup>
import { onMounted } from 'vue'
import { Settings, LogOut } from 'lucide-vue-next'
import { useSiteCheckStore } from '../stores/sitecheck'
import { SiteCheckService } from '../../bindings/sitecheck'

const store = useSiteCheckStore()

onMounted(() => {
  store.loadSettings()
})

const handleShowSettings = () => {
  SiteCheckService.ShowSettings()
}

const handleQuit = () => {
  SiteCheckService.Quit()
}
</script>

<template>
  <div class="tray-menu">
    <div class="tray-menu__list">
      <div v-for="target in store.settings.targets" :key="target.id" class="tray-menu__item tray-menu__item--disabled">
        <div class="tray-menu__target">
          <img v-if="target.iconUrl" :src="target.iconUrl" class="tray-menu__icon" alt="">
          <span class="tray-menu__name">{{ target.name }}</span>
        </div>
        <span class="tray-menu__latency" :class="{ 'tray-menu__latency--loading': store.loading }">
          {{ store.loading ? 'Testing...' : (store.report?.results?.find(r => r.id === target.id)?.latencyMs ? `${store.report.results.find(r => r.id === target.id).latencyMs}ms` : '--') }}
        </span>
      </div>
    </div>

    <div class="tray-menu__separator"></div>

    <div class="tray-menu__list">
      <button class="tray-menu__item tray-menu__button" @click="handleShowSettings">
        <div class="tray-menu__target">
          <Settings class="tray-menu__icon" />
          <span>Settings</span>
        </div>
      </button>
      <button class="tray-menu__item tray-menu__button" @click="handleQuit">
        <div class="tray-menu__target">
          <LogOut class="tray-menu__icon" />
          <span>Quit</span>
        </div>
      </button>
    </div>
  </div>
</template>

<style scoped>
.tray-menu {
  width: 100%;
  height: 100vh;
  padding: 12px;
  background: transparent;
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
  user-select: none;
  -webkit-user-select: none;
  -webkit-app-region: no-drag;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.tray-menu__list {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.tray-menu__item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 6px 10px;
  border-radius: 6px;
  font-size: 13px;
  border: none;
  background: transparent;
  width: 100%;
  text-align: left;
  color: inherit;
  cursor: default;
}

.tray-menu__item--disabled {
  opacity: 0.9;
}

.tray-menu__button {
  cursor: pointer;
}

.tray-menu__button:hover {
  background: #007aff;
  color: #fff;
}

.tray-menu__target {
  display: flex;
  align-items: center;
  gap: 10px;
}

.tray-menu__icon {
  width: 16px;
  height: 16px;
  object-fit: contain;
}

.tray-menu__name {
  font-weight: 500;
}

.tray-menu__latency {
  font-family: var(--font-mono);
  color: var(--muted);
}

.tray-menu__latency--loading {
  font-style: italic;
  font-size: 11px;
}

.tray-menu__separator {
  height: 1px;
  background: var(--separator);
  margin: 6px 4px;
}
</style>
