<script setup>
import { computed, nextTick, onBeforeUnmount, onMounted, ref } from "vue";
import { LogOut, Settings } from "lucide-vue-next";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { Spinner } from "@/components/ui/spinner";
import {
  Table,
  TableBody,
  TableCell,
  TableRow,
} from "@/components/ui/table";
import { countryFlagURL } from "@/lib/country-flags";
import { useSiteCheckStore } from "@/stores/sitecheck";
import { SiteCheckService } from "../../bindings/sitecheck";

const store = useSiteCheckStore();
const scrollContainer = ref(null);
const trayPageClass = "tray-menu-open";

const trayTargets = computed(() =>
  store.settings.targets.slice(0, 5).map((target) => {
    const result = store.report?.results?.find((entry) => entry.id === target.id);

    return {
      ...target,
      latency: store.loading
        ? "Testing..."
        : result?.status === "Available"
          ? `${result.latencyMs}ms`
          : result?.status === "Unavailable"
            ? "Unavailable"
          : "--",
    };
  }),
);

const myIPDisplay = computed(() => ({
  ip: store.myIPReport?.ip || "-",
  countryCode: normalizeCountryCode(store.myIPReport?.countryCode),
}));

const dnsCheckpoints = computed(() => {
  const checkpoints = store.dnsReport?.checkpoints;
  if (checkpoints?.length) {
    return checkpoints;
  }

  return Array.from({ length: 4 }, (_, index) => ({
    name: `check point #${index + 1}`,
    result: {
      isp: "--",
      ip: "--",
      country: "--",
    },
  }));
});

onMounted(() => {
  store.loadSettings();
  handleTrayActivated();
  document.documentElement.classList.add(trayPageClass);
  document.body.classList.add(trayPageClass);
  document.addEventListener("visibilitychange", handleVisibilityChange);
  window.addEventListener("focus", handleTrayActivated);
  scrollContainer.value?.addEventListener("wheel", handleWheel, { passive: false });
});

onBeforeUnmount(() => {
  scrollContainer.value?.removeEventListener("wheel", handleWheel);
  document.removeEventListener("visibilitychange", handleVisibilityChange);
  window.removeEventListener("focus", handleTrayActivated);
  document.documentElement.classList.remove(trayPageClass);
  document.body.classList.remove(trayPageClass);
});

function handleShowSettings() {
  SiteCheckService.ShowSettings();
}

function handleQuit() {
  SiteCheckService.Quit();
}

function fallbackLabel(name) {
  return (name || "?").slice(0, 2).toUpperCase();
}

function resetScrollTop() {
  nextTick(() => {
    if (scrollContainer.value) {
      scrollContainer.value.scrollTop = 0;
    }
  });
}

function handleVisibilityChange() {
  if (document.hidden) {
    resetScrollTop();
    return;
  }

  handleTrayActivated();
}

function handleTrayActivated() {
  store.markTrayRefresh();
}

function handleWheel(event) {
  const container = scrollContainer.value;
  if (!container) {
    return;
  }

  const maxScrollTop = container.scrollHeight - container.clientHeight;
  if (maxScrollTop <= 0) {
    event.preventDefault();
    return;
  }

  const atTop = container.scrollTop <= 0;
  const atBottom = container.scrollTop >= maxScrollTop - 1;

  if ((atTop && event.deltaY < 0) || (atBottom && event.deltaY > 0)) {
    event.preventDefault();
  }
}

function dnsValueClass(value) {
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

function normalizeCountryCode(value) {
  const code = String(value || "").trim().toLowerCase();
  if (!/^[a-z]{2}$/.test(code)) {
    return "";
  }
  return code;
}

function myIPFlagURL(code) {
  return code ? `https://flagcdn.com/w40/${code}.png` : "";
}
</script>

<template>
  <div class="tray-root flex h-screen bg-transparent">
    <div class="w-full overflow-hidden border border-border/70 bg-background/98 shadow-[0_12px_30px_-16px_rgba(15,23,42,0.22)] backdrop-blur-sm">
      <div ref="scrollContainer" class="tray-scroll h-full overflow-y-auto">
        <div class="flex flex-col gap-2 p-1.5">
          <section class="flex flex-col gap-1">
            <div class="flex items-center justify-between px-1 text-[11px] uppercase tracking-[0.08em] text-muted-foreground/80">
              <span>My IP</span>
              <Spinner v-if="store.myIPLoading" class="size-3.5" />
            </div>

            <Table>
              <TableBody>
                <TableRow class="border-none hover:bg-muted/20">
                  <TableCell class="py-1.5">
                    <div class="flex items-center gap-2">
                      <img
                        v-if="myIPFlagURL(myIPDisplay.countryCode)"
                        :src="myIPFlagURL(myIPDisplay.countryCode)"
                        :alt="myIPDisplay.countryCode.toUpperCase()"
                        class="h-3.5 w-5 rounded-[2px] border border-[#eaecf0] object-cover"
                      />
                      <span v-else class="text-[12px] text-muted-foreground">-</span>
                      <span class="text-[12px] font-medium tabular-nums text-foreground">
                        {{ myIPDisplay.ip }}
                      </span>
                    </div>
                  </TableCell>
                </TableRow>
              </TableBody>
            </Table>
          </section>

          <section class="flex flex-col gap-1">
            <div class="flex items-center justify-between px-1 text-[11px] uppercase tracking-[0.08em] text-muted-foreground/80">
              <span>Connectivity</span>
              <Spinner v-if="store.loading" class="size-3.5" />
            </div>

            <Table>
              <TableBody>
                <TableRow
                  v-for="target in trayTargets"
                  :key="target.id"
                  class="border-none hover:bg-muted/20"
                >
                  <TableCell class="py-1.5">
                    <div class="flex items-center gap-2">
                      <Avatar class="size-5 rounded-sm" shape="square">
                        <AvatarImage :src="target.iconUrl" :alt="target.name" />
                        <AvatarFallback>{{ fallbackLabel(target.name) }}</AvatarFallback>
                      </Avatar>
                      <span class="truncate text-[12px] font-medium">{{ target.name }}</span>
                    </div>
                  </TableCell>
                  <TableCell class="w-24 py-1.5 text-right">
                    <Badge variant="secondary" class="h-5 rounded-sm px-1.5 text-[10px] tabular-nums">
                      {{ target.latency }}
                    </Badge>
                  </TableCell>
                </TableRow>
              </TableBody>
            </Table>
          </section>

          <section class="flex flex-col gap-1">
            <div class="flex items-center justify-between px-1 text-[11px] uppercase tracking-[0.08em] text-muted-foreground/80">
              <span>DNS Leak Test</span>
              <Spinner v-if="store.dnsLoading" class="size-3.5" />
            </div>

            <div>
              <div
                v-for="checkpoint in dnsCheckpoints"
                :key="checkpoint.name"
                class="border-b px-2.5 py-2 last:border-b-0"
              >
                <div class="text-[12px] font-medium text-foreground">
                  {{ checkpoint.name }}
                </div>
                <div class="mt-1 flex flex-col gap-0.5 text-[11px]">
                  <div class="flex gap-1.5">
                    <span class="text-muted-foreground">ISP:</span>
                    <span class="truncate" :class="dnsValueClass(checkpoint.result.isp)">
                      {{ checkpoint.result.isp }}
                    </span>
                  </div>
                  <div class="flex gap-1.5">
                    <span class="text-muted-foreground">IP:</span>
                    <span class="flex min-w-0 items-center gap-1.5" :class="dnsValueClass(checkpoint.result.ip)">
                      <img
                        v-if="flagURL(checkpoint.result.country)"
                        :src="flagURL(checkpoint.result.country)"
                        :alt="checkpoint.result.country"
                        class="h-3 w-4 rounded-[2px] border border-[#eaecf0] object-cover"
                      />
                      <span class="truncate">
                        {{ checkpoint.result.ip }}
                      </span>
                    </span>
                  </div>
                </div>
              </div>
            </div>
          </section>

          <Separator />

          <div class="flex flex-col gap-1">
            <Button
              variant="ghost"
              size="sm"
              class="h-7 justify-start rounded-sm px-2 text-[12px] font-normal"
              @click="handleShowSettings"
            >
              <Settings data-icon="inline-start" />
              Settings
            </Button>
            <Button
              variant="ghost"
              size="sm"
              class="h-7 justify-start rounded-sm px-2 text-[12px] font-normal"
              @click="handleQuit"
            >
              <LogOut data-icon="inline-start" />
              Quit
            </Button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.tray-root {
  -webkit-app-region: no-drag;
  user-select: none;
  -webkit-user-select: none;
}

.tray-root :deep(*) {
  user-select: none;
  -webkit-user-select: none;
}

.tray-scroll {
  scrollbar-width: none;
  -ms-overflow-style: none;
  overscroll-behavior: none;
}

.tray-scroll::-webkit-scrollbar {
  display: none;
}

:global(html.tray-menu-open),
:global(body.tray-menu-open) {
  overscroll-behavior: none;
  overflow: hidden;
}

:global(body.tray-menu-open #app) {
  overflow: hidden;
}
</style>
