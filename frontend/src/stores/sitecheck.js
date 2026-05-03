import { computed, ref, shallowRef } from 'vue'
import { defineStore } from 'pinia'
import { Events } from '@wailsio/runtime'
import { SiteCheckService } from '../../bindings/sitecheck'

let toastTimer = 0

function defaultClientSettings() {
  return {
    intervalMinutes: 10,
    targets: [],
  }
}

function normalizeIntervalMinutes(value) {
  const interval = Number(value)
  if (!Number.isFinite(interval) || interval < 1) return 10
  return Math.min(Math.trunc(interval), 99)
}

function normalizeTargetURL(rawURL) {
  const value = String(rawURL || '').trim()
  if (!value) return null

  if (!/^https?:\/\//i.test(value)) {
    return null
  }

  try {
    const parsed = new URL(value)
    const host = parsed.hostname.replace(/^www\./, '')
    if (!host || !host.includes('.')) {
      return null
    }
    return value
  } catch {
    return null
  }
}

function buildCustomTarget(rawURL) {
  const normalizedURL = normalizeTargetURL(rawURL)
  if (!normalizedURL) return null

  const parsed = new URL(normalizedURL)
  const host = parsed.hostname.replace(/^www\./, '')

  return { id: '', name: host, url: normalizedURL, iconUrl: `https://favicon.im/${host}` }
}

function hasDuplicateTargetURL(targets, rawURL) {
  const trimmedURL = String(rawURL || '').trim()
  if (!trimmedURL) return false
  return targets.some((target) => String(target.url || '').trim() === trimmedURL)
}

export const useSiteCheckStore = defineStore('sitecheck', () => {
  const settings = ref(defaultClientSettings())
  const report = shallowRef(null)
  const loading = shallowRef(false)
  const saving = shallowRef(false)
  const message = shallowRef('')
  const messageTone = shallowRef('muted')
  const toastMessage = shallowRef('')
  const toastTone = shallowRef('muted')

  const canSave = computed(() => !saving.value && settings.value.targets.length >= 5)

  function setMessage(text, tone = 'muted') {
    message.value = text
    messageTone.value = tone
  }

  function showToast(text, tone = 'muted') {
    toastMessage.value = text
    toastTone.value = tone
    window.clearTimeout(toastTimer)
    toastTimer = window.setTimeout(() => {
      toastMessage.value = ''
      toastTone.value = 'muted'
    }, 2200)
  }

  async function loadSettings() {
    try {
      settings.value = await SiteCheckService.GetSettings()
      setMessage('')
    } catch (error) {
      setMessage(String(error), 'danger')
    }
  }

  function setIntervalMinutes(value) {
    const nextInterval = normalizeIntervalMinutes(value)
    settings.value = {
      ...settings.value,
      intervalMinutes: nextInterval,
    }
  }

  async function updateIntervalMinutes(value) {
    setIntervalMinutes(value)
    saving.value = true
    try {
      settings.value = await SiteCheckService.SaveSettings(settings.value)
      showToast(`Updated ${settings.value.intervalMinutes}m interval`, 'success')
    } catch (error) {
      showToast(String(error), 'danger')
    } finally {
      saving.value = false
    }
  }

  async function addTargetUrl(rawURL) {
    const nextTarget = buildCustomTarget(rawURL)
    if (!nextTarget) {
      showToast('Invalid site url', 'danger')
      return false
    }
    if (hasDuplicateTargetURL(settings.value.targets, rawURL)) {
      showToast('Site url already exists', 'danger')
      return false
    }

    saving.value = true
    try {
      settings.value = await SiteCheckService.SaveSettings({
        ...settings.value,
        targets: [nextTarget, ...settings.value.targets],
      })
      showToast(`Added ${nextTarget.name} to targets`, 'success')
      return true
    } catch (error) {
      showToast(String(error), 'danger')
      return false
    } finally {
      saving.value = false
    }
  }

  async function removeTarget(targetID) {
    const nextTargets = settings.value.targets.filter((target) => target.id !== targetID)
    if (nextTargets.length === settings.value.targets.length) {
      return false
    }

    saving.value = true
    try {
      settings.value = await SiteCheckService.SaveSettings({
        ...settings.value,
        targets: nextTargets,
      })
      showToast('Checkpoint removed', 'success')
      return true
    } catch (error) {
      showToast(String(error), 'danger')
      return false
    } finally {
      saving.value = false
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
    toastMessage,
    toastTone,
    canSave,
    loadSettings,
    setIntervalMinutes,
    updateIntervalMinutes,
    addTargetUrl,
    removeTarget,
    updateTarget,
    saveSettings,
    benchmark,
  }
})
