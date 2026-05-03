<script setup>
import { computed, ref, watch } from "vue";
import { CirclePlay, Clock3 } from "lucide-vue-next";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import { Spinner } from "@/components/ui/spinner";
import { countryFlagURL } from "@/lib/country-flags";
import { useSiteCheckStore } from "@/stores/sitecheck";

const store = useSiteCheckStore();
const intervalDraft = ref("1");

const checkpoints = computed(() => store.dnsReport?.checkpoints || []);
const hasResults = computed(() => checkpoints.value.length > 0);

const skeletonRows = [
  "check point #1",
  "check point #2",
  "check point #3",
  "check point #4",
];

watch(
  () => store.settings.dnsIntervalHours,
  (value) => {
    intervalDraft.value = String(value || 1);
  },
  { immediate: true },
);

function handleIntervalInput(event) {
  const raw = String(event.target.value ?? "");
  const digits = raw.replace(/\D/g, "");
  const interval = Number(digits);
  const normalized = !Number.isFinite(interval) || interval < 1
    ? 1
    : Math.min(Math.trunc(interval), 99);

  intervalDraft.value = String(normalized);
  store.updateDNSIntervalHours(normalized);
}

function valueClass(value) {
  if (value === "Timeout" || value === "Failed") {
    return "text-destructive";
  }
  if (value === "--") {
    return "text-muted-foreground";
  }
  return "text-foreground";
}

function flagURL(country) {
  return countryFlagURL(country);
}
</script>

<template>
  <section class="flex flex-col gap-4">
    <div class="flex flex-wrap items-center gap-3 border-b pb-3">
      <Button
        size="sm"
        class="h-8 rounded-sm px-3 text-[13px] shadow-none active:translate-y-px cursor-pointer"
        :disabled="store.dnsLoading"
        @click="store.benchmarkDNS"
      >
        <Spinner v-if="store.dnsLoading" data-icon="inline-start" />
        <CirclePlay v-else data-icon="inline-start" />
        {{ store.dnsLoading ? "Testing" : "Test" }}
      </Button>

      <div class="flex h-8 items-center gap-2 rounded-sm border bg-background px-2.5">
        <Clock3 class="size-3.5 text-muted-foreground" />
        <input
          type="number"
          min="1"
          max="99"
          :value="intervalDraft"
          class="h-full w-9 border-0 bg-transparent px-0 text-right text-[13px] tabular-nums shadow-none outline-none ring-0 focus-visible:outline-none focus-visible:ring-0 focus-visible:ring-offset-0"
          @input="handleIntervalInput"
        />
        <span class="text-[12px] text-muted-foreground">h</span>
      </div>
    </div>

    <div v-if="!hasResults" class="overflow-hidden rounded-sm border bg-card">
      <div
        v-for="label in skeletonRows"
        :key="label"
        class="flex flex-col gap-3 border-b px-3 py-3 last:border-b-0 sm:flex-row sm:items-start"
      >
        <div class="w-32 text-[13px] font-medium text-foreground">
          {{ label }}
        </div>

        <div class="grid min-w-0 flex-1 gap-3 sm:grid-cols-3">
          <div class="flex flex-col gap-1.5">
            <span class="text-[11px] uppercase tracking-[0.06em] text-muted-foreground">ISP</span>
            <Skeleton class="h-4 w-40 rounded-sm" />
          </div>
          <div class="flex flex-col gap-1.5">
            <span class="text-[11px] uppercase tracking-[0.06em] text-muted-foreground">IP</span>
            <Skeleton class="h-4 w-28 rounded-sm" />
          </div>
          <div class="flex flex-col gap-1.5">
            <span class="text-[11px] uppercase tracking-[0.06em] text-muted-foreground">Country</span>
            <Skeleton class="h-4 w-24 rounded-sm" />
          </div>
        </div>
      </div>
    </div>

    <div v-else class="overflow-hidden rounded-sm border bg-card">
      <div
        v-for="checkpoint in checkpoints"
        :key="checkpoint.name"
        class="flex flex-col gap-3 border-b px-3 py-3 last:border-b-0 sm:flex-row sm:items-start"
      >
        <div class="w-32 text-[13px] font-medium text-foreground">
          {{ checkpoint.name }}
        </div>

        <div class="grid min-w-0 flex-1 gap-3 sm:grid-cols-3">
          <div class="flex flex-col gap-1">
            <span class="text-[11px] uppercase tracking-[0.06em] text-muted-foreground">ISP</span>
            <span class="text-[13px] font-medium" :class="valueClass(checkpoint.result.isp)">
              {{ checkpoint.result.isp }}
            </span>
          </div>
          <div class="flex flex-col gap-1">
            <span class="text-[11px] uppercase tracking-[0.06em] text-muted-foreground">IP</span>
            <span class="text-[13px] font-medium" :class="valueClass(checkpoint.result.ip)">
              {{ checkpoint.result.ip }}
            </span>
          </div>
          <div class="flex flex-col gap-1">
            <span class="text-[11px] uppercase tracking-[0.06em] text-muted-foreground">Country</span>
            <span class="flex items-center gap-2 text-[13px] font-medium" :class="valueClass(checkpoint.result.country)">
              <img
                v-if="flagURL(checkpoint.result.country)"
                :src="flagURL(checkpoint.result.country)"
                :alt="checkpoint.result.country"
                class="h-3.5 w-5 rounded-[2px] border border-[#eaecf0] object-cover"
              />
              <span v-else>{{ checkpoint.result.country }}</span>
            </span>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>
