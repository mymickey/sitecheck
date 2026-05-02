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
        : result?.latencyMs
          ? `${result.latencyMs}ms`
          : "--",
    };
  }),
);

onMounted(() => {
  store.loadSettings();
  resetScrollTop();
  document.documentElement.classList.add(trayPageClass);
  document.body.classList.add(trayPageClass);
  document.addEventListener("visibilitychange", handleVisibilityChange);
  window.addEventListener("focus", resetScrollTop);
  scrollContainer.value?.addEventListener("wheel", handleWheel, { passive: false });
});

onBeforeUnmount(() => {
  scrollContainer.value?.removeEventListener("wheel", handleWheel);
  document.removeEventListener("visibilitychange", handleVisibilityChange);
  window.removeEventListener("focus", resetScrollTop);
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
  if (!document.hidden) {
    resetScrollTop();
  }
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
</script>

<template>
  <div class="tray-root flex h-screen bg-transparent">
    <div class="w-full overflow-hidden border border-border/70 bg-background/98 shadow-[0_12px_30px_-16px_rgba(15,23,42,0.22)] backdrop-blur-sm">
      <div ref="scrollContainer" class="tray-scroll h-full overflow-y-auto">
        <div class="flex flex-col gap-2 p-1.5">
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
                <TableCell class="w-20 py-1.5 text-right">
                  <Badge variant="secondary" class="h-5 rounded-sm px-1.5 text-[10px] tabular-nums">
                    <Spinner v-if="store.loading" data-icon="inline-start" />
                    {{ target.latency }}
                  </Badge>
                </TableCell>
              </TableRow>
            </TableBody>
          </Table>

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
