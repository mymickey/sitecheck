import { computed, ref, shallowRef } from 'vue'
import { defineStore } from 'pinia'
import { Events } from '@wailsio/runtime'
import { SiteCheckService } from '../../bindings/sitecheck'

const intervalStorageKey = 'sitecheck.intervalMinutes'

function defaultClientSettings() {
  return {
    intervalMinutes: 10,
    targets: [],
  }
}

function readStoredIntervalMinutes() {
  const raw = window.localStorage.getItem(intervalStorageKey)
  const interval = Number(raw)
  return Number.isFinite(interval) && interval > 0 ? interval : 10
}

function writeStoredIntervalMinutes(value) {
  window.localStorage.setItem(intervalStorageKey, String(value))
}

function normalizeIntervalMinutes(value) {
  const interval = Number(value)
  if (!Number.isFinite(interval) || interval < 1) return 10
  return Math.min(Math.trunc(interval), 99)
}

export const useSiteCheckStore = defineStore('sitecheck', () => {
  const settings = ref(defaultClientSettings())
  const report = shallowRef(null)
  const loading = shallowRef(false)
  const saving = shallowRef(false)
  const message = shallowRef('')
  const messageTone = shallowRef('muted')

  const canSave = computed(() => !saving.value && settings.value.targets.length === 5)

  function setMessage(text, tone = 'muted') {
    message.value = text
    messageTone.value = tone
  }

  async function loadSettings() {
    try {
      const nextSettings = await SiteCheckService.GetSettings()
      settings.value = {
        ...nextSettings,
        intervalMinutes: readStoredIntervalMinutes(),
      }
      setMessage('')
    } catch (error) {
      setMessage(String(error), 'danger')
    }
  }

  function setIntervalMinutes(value) {
    const nextInterval = normalizeIntervalMinutes(value)
    writeStoredIntervalMinutes(nextInterval)
    settings.value = {
      ...settings.value,
      intervalMinutes: nextInterval,
    }
  }

  function updateTarget(index, patch) {
    settings.value = {
      ...settings.value,
      targets: settings.value.targets.map((target, targetIndex) => (
        targetIndex === index ? { ...target, ...patch } : target
      )),
    }
  }

  async function saveSettings() {
    if (!canSave.value) return
    saving.value = true
    try {
      settings.value = await SiteCheckService.SaveSettings(settings.value)
      setMessage('Saved', 'success')
    } catch (error) {
      setMessage(String(error), 'danger')
    } finally {
      saving.value = false
    }
  }

  async function benchmark() {
    loading.value = true
    try {
      report.value = await SiteCheckService.Benchmark()
      setMessage('Benchmark complete', 'success')
    } catch (error) {
      setMessage(String(error), 'danger')
    } finally {
      loading.value = false
    }
  }

  Events.On('benchmark-finished', (event) => {
    report.value = event.data
    loading.value = false
  })

  Events.On('settings-updated', (event) => {
    settings.value = event.data
  })

  return {
    settings,
    report,
    loading,
    saving,
    message,
    messageTone,
    canSave,
    loadSettings,
    setIntervalMinutes,
    updateTarget,
    saveSettings,
    benchmark,
  }
})
