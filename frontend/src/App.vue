<script setup>
import { computed, onMounted, shallowRef } from "vue";
import { LogOut } from "lucide-vue-next";
import AboutPanel from "@/components/AboutPanel.vue";
import ConnectivitySettings from "@/components/ConnectivitySettings.vue";
import TrayMenu from "@/components/TrayMenu.vue";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Separator } from "@/components/ui/separator";
import { useSiteCheckStore } from "@/stores/sitecheck";
import { SiteCheckService } from "../bindings/sitecheck";

const isMenuMode = window.location.search.includes("mode=menu");
const activeTab = shallowRef("connectivity");

const tabs = [
  { id: "connectivity", label: "Connectivity" },
  { id: "about", label: "Help & Information" },
];

const activeComponent = computed(() =>
  activeTab.value === "about" ? AboutPanel : ConnectivitySettings,
);

onMounted(() => {
  useSiteCheckStore().loadSettings();
});

function handleQuit() {
  SiteCheckService.Quit();
}
</script>

<template>
  <div :class="isMenuMode ? 'h-full bg-transparent' : 'h-full bg-background'">
    <TrayMenu v-if="isMenuMode" />

    <div v-else class="app-shell flex h-screen flex-col">
      <header class="window-header shrink-0"></header>

      <main class="grid min-h-0 flex-1 grid-cols-[244px_minmax(0,1fr)] gap-0 p-2">
        <aside class="min-h-0">
          <Card class="flex h-full flex-col rounded-none border-0 bg-sidebar shadow-none">
            <CardHeader class="gap-4 p-5 pb-4">
              <div class="flex items-center gap-3">
                <Avatar class="size-8 border border-border/70 bg-background" shape="circle">
                  <AvatarImage src="/logo.svg" alt="SiteCheck" />
                  <AvatarFallback>SC</AvatarFallback>
                </Avatar>
                <div class="flex min-w-0 flex-col gap-1">
                  <CardTitle class="text-[15px] font-semibold">SiteCheck</CardTitle>
                  <CardDescription class="text-xs">Connectivity monitor</CardDescription>
                </div>
              </div>
            </CardHeader>

            <CardContent class="flex min-h-0 flex-1 flex-col gap-6 px-3 pb-3 pt-0">
              <div class="flex flex-col gap-1">
                <div class="px-2 pb-1 text-xs text-muted-foreground">Settings</div>
                <Button
                  v-for="tab in tabs"
                  :key="tab.id"
                  :variant="activeTab === tab.id ? 'secondary' : 'ghost'"
                  size="sm"
                  class="h-9 justify-start rounded-md font-normal"
                  @click="activeTab = tab.id"
                >
                  {{ tab.label }}
                </Button>
              </div>

              <div class="mt-auto flex flex-col gap-2">
                <Separator />
                <div class="px-2 pt-1 text-xs text-muted-foreground">App</div>
                <Button variant="ghost" size="sm" class="h-9 justify-start rounded-md font-normal" @click="handleQuit">
                  <LogOut data-icon="inline-start" />
                  Log out
                </Button>
              </div>
            </CardContent>
          </Card>
        </aside>

        <section class="flex min-h-0 flex-col">
          <Card class="flex min-h-0 flex-1 flex-col rounded-none border-0 bg-transparent shadow-none">
            <ScrollArea class="flex-1">
              <div class="p-5 pt-4">
                <component :is="activeComponent" />
              </div>
            </ScrollArea>
          </Card>
        </section>
      </main>
    </div>
  </div>
</template>

<style scoped>
.window-header {
  height: 18px;
  -webkit-app-region: drag;
}
</style>
