<script setup>
import { BookOpen, Github, Globe, Wrench } from "lucide-vue-next";
import { Browser } from "@wailsio/runtime";
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
  { label: "Platform", value: "macOS only", icon: BookOpen },
];

const links = [
  {
    label: "GitHub",
    value: "github.com/mymickey/sitecheck",
    href: "https://github.com/mymickey/sitecheck",
    icon: Github,
  },
  {
    label: "Website",
    value: "sitecheck.top",
    href: "https://sitecheck.top",
    icon: Globe,
  },
];

function openExternalURL(url) {
  void Browser.OpenURL(url);
}
</script>

<template>
  <section class="grid gap-4">
    <Card>
      <CardHeader class="gap-4 pb-3">
        <div class="flex items-center gap-4">
          <Avatar size="base" shape="square">
            <AvatarImage src="/sitecheck-app.png" alt="SiteCheck" />
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

        <div class="grid gap-3 sm:grid-cols-2">
          <button
            v-for="link in links"
            :key="link.label"
            type="button"
            class="rounded-lg border bg-muted/30 p-3 text-left transition-colors hover:bg-muted/50"
            @click="openExternalURL(link.href)"
          >
            <div class="flex items-center gap-2 text-xs text-muted-foreground">
              <component :is="link.icon" class="size-4" />
              {{ link.label }}
            </div>
            <div class="mt-2 text-sm font-semibold">{{ link.value }}</div>
          </button>
        </div>
      </CardContent>
    </Card>
  </section>
</template>
