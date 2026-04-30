<script setup>
import { computed, nextTick, onBeforeUnmount, onMounted, ref } from "vue";
import { LogOut, Settings } from "lucide-vue-next";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
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

const trayTargets = computed(() =>
  store.settings.targets.map((target) => {
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
  document.addEventListener("visibilitychange", handleVisibilityChange);
  window.addEventListener("focus", resetScrollTop);
});

onBeforeUnmount(() => {
  document.removeEventListener("visibilitychange", handleVisibilityChange);
  window.removeEventListener("focus", resetScrollTop);
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
</script>

<template>
  <div class="tray-root flex h-screen bg-transparent">
    <Card class="w-full overflow-hidden rounded-none border-0 shadow-none">
      <div ref="scrollContainer" class="tray-scroll h-full overflow-y-auto">
        <CardContent class="flex flex-col gap-2 p-1.5">
          <Table>
            <TableBody>
              <TableRow
                v-for="target in trayTargets"
                :key="target.id"
                class="hover:bg-muted/20 border-none"
              >
                <TableCell class="py-2">
                  <div class="flex items-center gap-2.5">
                    <Avatar class="size-5 rounded-sm" shape="square">
                      <AvatarImage :src="target.iconUrl" :alt="target.name" />
                      <AvatarFallback>{{ fallbackLabel(target.name) }}</AvatarFallback>
                    </Avatar>
                    <span class="truncate text-[13px] font-medium">{{ target.name }}</span>
                  </div>
                </TableCell>
                <TableCell class="w-24 py-2 text-right">
                  <Badge variant="outline" class="rounded-full text-[11px]">
                    <Spinner v-if="store.loading" data-icon="inline-start" />
                    {{ target.latency }}
                  </Badge>
                </TableCell>
              </TableRow>
            </TableBody>
          </Table>

          <Separator />

          <div class="flex flex-col gap-1">
            <Button variant="ghost" size="sm" class="h-8 justify-start rounded-md px-2.5 font-normal" @click="handleShowSettings">
              <Settings data-icon="inline-start" />
              Settings
            </Button>
            <Button variant="ghost" size="sm" class="h-8 justify-start rounded-md px-2.5 font-normal" @click="handleQuit">
              <LogOut data-icon="inline-start" />
              Quit
            </Button>
          </div>
        </CardContent>
      </div>
    </Card>
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
}

.tray-scroll::-webkit-scrollbar {
  display: none;
}
</style>
