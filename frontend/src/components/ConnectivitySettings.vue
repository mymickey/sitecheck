<script setup>
import { computed } from 'vue'
import { Play, RotateCw, Save } from 'lucide-vue-next'
import { useSiteCheckStore } from '../stores/sitecheck'

const store = useSiteCheckStore()

const resultsById = computed(() => {
  const entries = store.report?.results || []
  return Object.fromEntries(entries.map((result) => [result.id, result]))
})

const summaryLabel = computed(() => {
  if (!store.report?.summary?.hasResults) return { fast: '--', slow: '--' }
  return {
    fast: `${store.report.summary.fastestMs}ms`,
    slow: `${store.report.summary.slowestMs}ms`,
  }
})

function targetResult(target) {
  return resultsById.value[target.id] || null
}

function resultTone(result) {
  if (!result) return 'idle'
  if (result.status === 'Available') return result.latencyMs < 200 ? 'success' : 'warning'
  return 'danger'
}

function resultText(result) {
  if (!result) return 'Waiting'
  if (result.status === 'Available') return `${result.latencyMs}ms`
  return result.status
}
</script>

<template>
  <form class="panel connectivity" @submit.prevent="store.saveSettings">
    <header class="panel__header">
      <div>
        <h1 class="panel__title">Connectivity</h1>
        <p class="panel__subtitle">Network targets and background benchmark interval.</p>
      </div>
      <div class="toolbar">
        <button
          class="button button--secondary"
          type="button"
          :disabled="store.loading"
          @click="store.benchmark"
        >
          <component :is="store.loading ? RotateCw : Play" class="button__icon" aria-hidden="true" />
          <span>Benchmark</span>
        </button>
        <button class="button button--primary" type="submit" :disabled="store.saving">
          <Save class="button__icon" aria-hidden="true" />
          <span>Save</span>
        </button>
      </div>
    </header>

    <div class="metric-strip" aria-label="Latest benchmark summary">
      <div class="metric">
        <span class="metric__label">Fastest</span>
        <strong>{{ summaryLabel.fast }}</strong>
      </div>
      <div class="metric">
        <span class="metric__label">Slowest</span>
        <strong>{{ summaryLabel.slow }}</strong>
      </div>
      <label class="interval-field">
        <span>Interval</span>
        <input
          class="input input--secondary interval-field__input"
          type="number"
          min="1"
          max="1440"
          :value="store.settings.intervalMinutes"
          @input="store.setIntervalMinutes($event.target.value)"
        >
        <span>min</span>
      </label>
    </div>

    <div class="target-list">
      <article v-for="(target, index) in store.settings.targets" :key="target.id" class="target-row">
        <img class="target-row__icon" :src="target.iconUrl" alt="" aria-hidden="true">
        <div class="target-row__fields">
          <label class="field">
            <span>Name</span>
            <input
              class="input input--secondary input--full-width"
              :value="target.name"
              @input="store.updateTarget(index, { name: $event.target.value })"
            >
          </label>
          <label class="field field--url">
            <span>URL</span>
            <input
              class="input input--secondary input--full-width"
              :value="target.url"
              @input="store.updateTarget(index, { url: $event.target.value })"
            >
          </label>
        </div>
        <span class="status-chip" :data-tone="resultTone(targetResult(target))">
          {{ resultText(targetResult(target)) }}
        </span>
      </article>
    </div>

    <footer class="panel__footer">
      <span class="message" :data-tone="store.messageTone">{{ store.message }}</span>
    </footer>
  </form>
</template>
