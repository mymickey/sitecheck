<script setup>
import { Clock3, Gauge, ShieldCheck, TimerReset } from "lucide-vue-next";
import { Badge } from "@/components/ui/badge";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Progress } from "@/components/ui/progress";
import { Separator } from "@/components/ui/separator";

defineProps({
  summaryCards: { type: Array, required: true },
  intervalMinutes: { type: Number, required: true },
  averageLabel: { type: String, required: true },
  availableTargets: { type: Number, required: true },
  targetRows: { type: Array, required: true },
});
</script>

<template>
  <div class="flex flex-col gap-6">
    <div class="grid gap-4 xl:grid-cols-2">
      <Card v-for="card in summaryCards" :key="card.label">
        <CardContent class="flex items-start justify-between p-6">
          <div class="flex flex-col gap-3">
            <span class="text-sm text-muted-foreground">{{ card.label }}</span>
            <span class="text-4xl font-semibold tracking-tight">{{ card.value }}</span>
            <span class="text-sm font-medium">{{ card.delta }}</span>
          </div>
          <Badge variant="outline" class="rounded-full px-3 py-1 text-sm">{{ card.label }}</Badge>
        </CardContent>
      </Card>
    </div>

    <Card>
      <CardHeader class="gap-2 pb-3">
        <div class="flex items-center justify-between gap-3">
          <div class="flex flex-col gap-1">
            <CardTitle class="text-sm">Connectivity Snapshot</CardTitle>
            <CardDescription class="text-xs">
              Current route health and target latency distribution.
            </CardDescription>
          </div>
          <Badge variant="secondary">
            <Clock3 data-icon="inline-start" />
            {{ intervalMinutes }} min
          </Badge>
        </div>
      </CardHeader>
      <CardContent class="flex flex-col gap-4">
        <div class="grid gap-3 md:grid-cols-2">
          <div class="rounded-lg border bg-background p-3">
            <div class="flex items-center gap-2 text-xs text-muted-foreground">
              <Gauge class="size-4" />
              Average latency
            </div>
            <div class="mt-2 text-lg font-semibold tracking-tight">{{ averageLabel }}</div>
          </div>

          <div class="rounded-lg border bg-background p-3">
            <div class="flex items-center gap-2 text-xs text-muted-foreground">
              <ShieldCheck class="size-4" />
              Reachable targets
            </div>
            <div class="mt-2 text-lg font-semibold tracking-tight">
              {{ availableTargets }}/5
            </div>
          </div>
        </div>

        <Separator />

        <div class="flex flex-col gap-3">
          <div
            v-for="target in targetRows"
            :key="`${target.id}-status`"
            class="flex items-center gap-3 rounded-lg border bg-background px-3 py-2.5"
          >
            <div class="flex min-w-0 flex-1 items-center gap-3">
              <div class="flex min-w-0 flex-1 flex-col gap-1">
                <div class="flex items-center gap-2">
                  <span class="truncate text-sm font-medium">{{ target.name }}</span>
                  <span class="text-xs text-muted-foreground">{{ target.statusText }}</span>
                </div>
                <Progress :model-value="target.progress" class="h-1.5" />
              </div>
            </div>
            <Badge :variant="target.statusVariant">
              <TimerReset v-if="target.statusText === 'Waiting'" data-icon="inline-start" />
              {{ target.statusText }}
            </Badge>
          </div>
        </div>
      </CardContent>
    </Card>
  </div>
</template>
