<script setup>
import { ref, watch } from "vue";
import { CirclePlay, Clock3 } from "lucide-vue-next";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Spinner } from "@/components/ui/spinner";

const props = defineProps({
  targetRows: { type: Array, required: true },
  intervalMinutes: { type: Number, required: true },
  availableTargets: { type: Number, required: true },
  averageLabel: { type: String, required: true },
  loading: { type: Boolean, required: true },
});

const emit = defineEmits(["benchmark", "update-interval"]);

function fallbackLabel(name) {
  return (name || "?").slice(0, 2).toUpperCase();
}

const intervalDraft = ref("10");

watch(
  () => props.intervalMinutes,
  (value) => {
    intervalDraft.value = String(value || 10);
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
  emit("update-interval", normalized);
}
</script>

<template>
  <section class="flex flex-col gap-4">
    <div class="flex flex-wrap items-center justify-between gap-3 border-b pb-3">
      <div class="flex items-center gap-2">
        <Button
          size="sm"
          class="h-8 rounded-sm px-3 text-[13px] shadow-none active:translate-y-px cursor-pointer"
          :disabled="loading"
          @click="emit('benchmark')"
        >
          <Spinner v-if="loading" data-icon="inline-start" />
          <CirclePlay v-else data-icon="inline-start" />
          {{ loading ? "Benchmarking" : "Benchmark" }}
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
          <span class="text-[12px] text-muted-foreground">m</span>
        </div>
      </div>

      <div class="flex items-center gap-2">
        <Badge variant="secondary" class="h-6 rounded-sm px-2 text-[11px] font-medium">
          {{ availableTargets }}/5 reachable
        </Badge>
        <Badge variant="secondary" class="h-6 rounded-sm px-2 text-[11px] font-medium">
          {{ averageLabel }}
        </Badge>
      </div>
    </div>

    <div class="overflow-hidden rounded-sm border bg-card">
      <div
        v-for="target in targetRows"
        :key="target.id"
        class="flex items-center gap-3 border-b px-3 py-2.5 last:border-b-0"
      >
        <Avatar class="size-5 rounded-sm" shape="square">
          <AvatarImage :src="target.iconUrl" :alt="target.name" />
          <AvatarFallback>{{ fallbackLabel(target.name) }}</AvatarFallback>
        </Avatar>

        <div class="flex min-w-0 flex-1 items-center gap-3">
          <div class="w-36 truncate text-[13px] font-medium text-foreground">
            {{ target.name }}
          </div>

          <div class="min-w-0 flex-1 truncate text-[12px] text-muted-foreground">
            {{ target.url }}
          </div>
        </div>

        <Badge class="h-6 rounded-sm px-2 text-[11px] tabular-nums" :variant="target.statusVariant">
          {{ target.statusText }}
        </Badge>
      </div>
    </div>
  </section>
</template>

<style scoped>
</style>
