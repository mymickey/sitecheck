<script setup>
import { BookOpen, MonitorSmartphone, Radar, Wrench } from "lucide-vue-next";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

const facts = [
  { label: "Version", value: "0.1.0", icon: Wrench },
  { label: "Runtime", value: "Wails v3", icon: MonitorSmartphone },
  { label: "Probe Method", value: "GET + no-store", icon: Radar },
  { label: "Platform", value: "macOS only", icon: BookOpen },
];
</script>

<template>
  <section class="grid gap-4">
    <Card>
      <CardHeader class="gap-4 pb-3">
        <div class="flex items-center gap-4">
          <Avatar size="base" shape="square">
            <AvatarImage src="/logo.svg" alt="SiteCheck" />
            <AvatarFallback>SC</AvatarFallback>
          </Avatar>
          <div class="flex flex-col gap-1">
            <CardTitle class="text-sm">SiteCheck</CardTitle>
            <CardDescription class="text-xs">
              macOS menubar connectivity monitor built around the IPCheck benchmark flow.
            </CardDescription>
          </div>
        </div>
      </CardHeader>
      <CardContent class="flex flex-col gap-4">
        <div class="grid gap-3 sm:grid-cols-2">
          <div
            v-for="fact in facts"
            :key="fact.label"
            class="rounded-lg border bg-muted/30 p-3"
          >
            <div class="flex items-center gap-2 text-xs text-muted-foreground">
              <component :is="fact.icon" class="size-4" />
              {{ fact.label }}
            </div>
            <div class="mt-2 text-sm font-semibold">{{ fact.value }}</div>
          </div>
        </div>
      </CardContent>
    </Card>
  </section>
</template>
