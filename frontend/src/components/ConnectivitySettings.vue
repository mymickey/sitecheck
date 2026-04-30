<script setup>
import { computed } from "vue";
import ConnectivityTargetsTable from "@/components/connectivity/ConnectivityTargetsTable.vue";
import { useSiteCheckStore } from "@/stores/sitecheck";

const store = useSiteCheckStore();

const resultsById = computed(() => {
  const entries = store.report?.results || [];

  return Object.fromEntries(entries.map((result) => [result.id, result]));
});

const targetRows = computed(() =>
  store.settings.targets.map((target, index) => {
    const result = resultsById.value[target.id] || null;

    return {
      ...target,
      index,
      host: hostLabel(target),
      statusText: resultText(result),
      statusVariant: statusVariant(result),
      detail: detailText(result),
      progress: progressValue(result),
    };
  }),
);

const availableTargets = computed(
  () =>
    targetRows.value.filter((target) => target.statusText.endsWith("ms")).length,
);

const averageLabel = computed(() => {
  const values = targetRows.value
    .map((target) => latencyValue(target.statusText))
    .filter(Boolean);

  if (!values.length) return "No benchmark yet";

  const average = Math.round(
    values.reduce((total, value) => total + value, 0) / values.length,
  );

  return `${average}ms average`;
});

function latencyValue(statusText) {
  return statusText.endsWith("ms") ? Number.parseInt(statusText, 10) : 0;
}

function statusVariant(result) {
  if (!result) return "outline";
  if (result.status === "Available") {
    return result.latencyMs < 200 ? "default" : "secondary";
  }

  return "destructive";
}

function resultText(result) {
  if (!result) return "Waiting";
  if (result.status === "Available") return `${result.latencyMs}ms`;

  return result.status;
}

function detailText(result) {
  if (!result) return "Awaiting the first benchmark cycle.";
  if (result.status === "Available") return "Probe reachable from the current route.";

  return "Connection, DNS, or timeout failure in the latest probe.";
}

function progressValue(result) {
  if (!result || result.status !== "Available") return 8;

  const all = (store.report?.results || [])
    .filter((entry) => entry.status === "Available")
    .map((entry) => entry.latencyMs);
  const maximum = Math.max(...all, 1);

  return Math.max(Math.round((result.latencyMs / maximum) * 100), 12);
}

function hostLabel(target) {
  try {
    return new URL(target.url).hostname;
  } catch {
    return target.url;
  }
}
</script>

<template>
  <section class="flex flex-col gap-6">
    <ConnectivityTargetsTable
      :target-rows="targetRows"
      :interval-minutes="store.settings.intervalMinutes"
      :available-targets="availableTargets"
      :average-label="averageLabel"
      :loading="store.loading"
      @benchmark="store.benchmark"
      @update-interval="store.updateIntervalMinutes"
    />
  </section>
</template>
