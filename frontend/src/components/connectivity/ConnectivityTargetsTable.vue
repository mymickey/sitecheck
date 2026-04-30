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
  intervalDraft.value = event.target.value;
}

function handleIntervalBlur() {
  const interval = Number(intervalDraft.value);
  if (!Number.isFinite(interval)) {
    intervalDraft.value = "10";
    emit("update-interval", 10);
    return;
  }

  if (interval < 1) {
    intervalDraft.value = "1";
    emit("update-interval", 1);
    return;
  }

  if (interval > 99) {
    intervalDraft.value = "99";
    emit("update-interval", 99);
    return;
  }

  const normalized = Math.trunc(interval);
  intervalDraft.value = String(normalized);
  emit("update-interval", normalized);
}
</script>

<template>
  <section class="flex flex-col gap-3">
    <div class="flex flex-wrap items-center justify-between gap-3">
      <div class="flex items-center gap-2">
        <Button size="sm" class="rounded-md cursor-pointer" :disabled="loading" @click="emit('benchmark')">
          <Spinner v-if="loading" data-icon="inline-start" />
          <CirclePlay v-else data-icon="inline-start" />
          {{ loading ? "Benchmarking" : "Benchmark" }}
        </Button>
        <div class="flex h-9 items-center gap-2 rounded-md border bg-background px-3">
          <Clock3 class="size-4 text-muted-foreground" />
          <input
            type="number"
            min="1"
            max="99"
            :value="intervalDraft"
            class="h-full w-10 border-0 bg-transparent px-0 text-right text-sm shadow-none outline-none ring-0 focus-visible:outline-none focus-visible:ring-0 focus-visible:ring-offset-0"
            @input="handleIntervalInput"
            @blur="handleIntervalBlur"
          />
          <span class="text-sm text-muted-foreground">m</span>
        </div>
      </div>

      <div class="flex items-center gap-2">
        <Badge variant="secondary">{{ availableTargets }}/5 reachable</Badge>
        <Badge variant="secondary">{{ averageLabel }}</Badge>
      </div>
    </div>

    <div class="flex flex-col gap-3">
      <div
        v-for="target in targetRows"
        :key="target.id"
        class="flex items-center gap-3 rounded-lg border bg-background px-3 py-2.5"
      >
        <Avatar class="size-5 rounded-sm" shape="square">
          <AvatarImage :src="target.iconUrl" :alt="target.name" />
          <AvatarFallback>{{ fallbackLabel(target.name) }}</AvatarFallback>
        </Avatar>

        <div class="flex min-w-0 flex-1 items-center gap-3">
          <div class="w-40 truncate text-sm font-medium text-foreground">
            {{ target.name }}
          </div>

          <div class="min-w-0 flex-1 truncate text-sm text-muted-foreground">
            {{ target.url }}
          </div>
        </div>

        <Badge class="rounded-full" :variant="target.statusVariant">{{ target.statusText }}</Badge>
      </div>
    </div>
  </section>
</template>

<style scoped>
</style>
