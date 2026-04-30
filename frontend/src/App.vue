<script setup>
import { computed, onMounted, shallowRef } from "vue";
import AboutPanel from "@/components/AboutPanel.vue";
import ConnectivitySettings from "@/components/ConnectivitySettings.vue";
import TrayMenu from "@/components/TrayMenu.vue";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import { CardDescription, CardTitle } from "@/components/ui/card";
import { useSiteCheckStore } from "@/stores/sitecheck";

const isMenuMode = window.location.search.includes("mode=menu");
const activeTab = shallowRef("connectivity");
const store = useSiteCheckStore();

const tabs = [
  { id: "connectivity", label: "Connectivity" },
  { id: "about", label: "Help & Information" },
];

const activeComponent = computed(() =>
  activeTab.value === "about" ? AboutPanel : ConnectivitySettings,
);

const toastClass = computed(() => {
  if (store.toastTone === "danger") {
    return "border-destructive/20 bg-destructive text-destructive-foreground";
  }

  return "border-border bg-foreground text-background";
});

onMounted(() => {
  store.loadSettings();
});
</script>

<template>
  <div :class="isMenuMode ? 'h-full bg-transparent' : 'h-full bg-background'">
    <TrayMenu v-if="isMenuMode" />

    <div v-else class="app-shell flex h-screen flex-col">
      <header class="window-header shrink-0"></header>

      <main class="grid min-h-0 flex-1 grid-cols-[220px_minmax(0,1fr)]">
        <aside class="min-h-0 border-r bg-sidebar/65">
          <div class="flex h-full flex-col px-4 py-4">
            <div class="flex items-center gap-3 px-2 pb-4">
              <Avatar class="size-8 border border-border/70 bg-background" shape="circle">
                <AvatarImage src="/logo.svg" alt="SiteCheck" />
                <AvatarFallback>SC</AvatarFallback>
              </Avatar>
              <div class="flex min-w-0 flex-col gap-0.5">
                <CardTitle class="text-[14px] font-semibold tracking-tight">SiteCheck</CardTitle>
                <CardDescription class="text-[11px]">Connectivity monitor</CardDescription>
              </div>
            </div>

            <div class="flex flex-col gap-1 border-t pt-4">
              <div class="px-2 pb-1 text-[11px] uppercase tracking-[0.08em] text-muted-foreground/80">
                Settings
              </div>
              <Button
                v-for="tab in tabs"
                :key="tab.id"
                :variant="activeTab === tab.id ? 'secondary' : 'ghost'"
                size="sm"
                class="h-8 justify-start rounded-sm px-2.5 text-[13px] font-normal shadow-none"
                @click="activeTab = tab.id"
              >
                {{ tab.label }}
              </Button>
            </div>
          </div>
        </aside>

        <section class="flex min-h-0 flex-col bg-background">
          <div class="flex-1 overflow-y-auto">
            <div class="p-5 pt-4">
              <component :is="activeComponent" />
            </div>
          </div>
        </section>
      </main>

      <div
        v-if="store.toastMessage"
        class="pointer-events-none fixed right-5 top-8 z-50"
      >
        <div
          :class="toastClass"
          class="rounded-sm border px-3 py-2 text-[12px] font-medium shadow-[0_14px_32px_-18px_rgba(15,23,42,0.35)]"
        >
          {{ store.toastMessage }}
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.window-header {
  height: 18px;
  -webkit-app-region: drag;
}
</style>
