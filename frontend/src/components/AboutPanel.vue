<script setup>
import { BookOpen, ExternalLink, MonitorSmartphone, Radar, Wrench } from "lucide-vue-next";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";

const links = [
  {
    label: "Wails v3 Alpha",
    meta: "Desktop runtime, native menu integration, and accessory app lifecycle.",
    href: "https://v3alpha.wails.io/",
  },
  {
    label: "IPCheck.ing",
    meta: "Reference implementation behind the connectivity benchmark behavior.",
    href: "https://ipcheck.ing/",
  },
];

const facts = [
  { label: "Version", value: "0.1.0", icon: Wrench },
  { label: "Runtime", value: "Wails v3", icon: MonitorSmartphone },
  { label: "Probe Method", value: "GET + no-store", icon: Radar },
  { label: "Platform", value: "macOS only", icon: BookOpen },
];
</script>

<template>
  <section class="grid gap-4 xl:grid-cols-[minmax(0,1.1fr)_minmax(320px,0.9fr)]">
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

        <Separator />

        <div class="flex flex-col gap-3 text-sm text-muted-foreground">
          <p>
            SiteCheck focuses on the connectivity portion of IPCheck: practical
            route reachability to your chosen sites from the current network path.
          </p>
          <p>
            Each target is requested directly, any HTTP response counts as reachable,
            and DNS, connection, or timeout failures are shown as unavailable.
          </p>
        </div>
      </CardContent>
      <CardFooter class="pt-0">
        <Badge variant="secondary">
          Paid-tool polish means staying quiet, legible, and native-looking while the data is live.
        </Badge>
      </CardFooter>
    </Card>

    <Card>
      <CardHeader class="pb-3">
        <CardTitle class="text-sm">Resources</CardTitle>
        <CardDescription class="text-xs">
          Primary references behind the current implementation.
        </CardDescription>
      </CardHeader>
      <CardContent class="flex flex-col gap-3">
        <div
          v-for="link in links"
          :key="link.href"
          class="flex flex-col gap-3 rounded-lg border p-3"
        >
          <div class="flex flex-col gap-1">
            <div class="text-sm font-medium">{{ link.label }}</div>
            <p class="text-xs text-muted-foreground">{{ link.meta }}</p>
          </div>
          <Button as="a" size="sm" :data-wml-openURL="link.href" variant="outline" class="w-fit">
            <ExternalLink data-icon="inline-start" />
            Open
          </Button>
        </div>
      </CardContent>
    </Card>
  </section>
</template>
