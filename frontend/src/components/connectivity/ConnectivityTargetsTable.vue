<script setup>
import { computed, onBeforeUnmount, onMounted, ref, watch } from "vue";
import { CirclePlay, Clock3, Plus } from "lucide-vue-next";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  DialogTitle,
} from "reka-ui";
import { Modal } from "@/components/ui/modal";
import { Spinner } from "@/components/ui/spinner";

const props = defineProps({
  targetRows: { type: Array, required: true },
  intervalMinutes: { type: Number, required: true },
  availableTargets: { type: Number, required: true },
  averageLabel: { type: String, required: true },
  loading: { type: Boolean, required: true },
  saving: { type: Boolean, required: true },
  addTargetUrl: { type: Function, required: true },
  removeTarget: { type: Function, required: false, default: null },
});

const emit = defineEmits(["benchmark", "update-interval"]);

function fallbackLabel(name) {
  return (name || "?").slice(0, 2).toUpperCase();
}

const intervalDraft = ref("10");
const modalOpen = ref(false);
const siteURL = ref("");
const siteURLTouched = ref(false);
const activeDeleteTargetID = ref("");
const rootElement = ref(null);

watch(
  () => props.intervalMinutes,
  (value) => {
    intervalDraft.value = String(value || 10);
  },
  { immediate: true },
);

function handleIntervalInput(event) {
  const raw = String(event.target.value ?? "");
  const digits = raw.replace(/\D/g, "");
  const interval = Number(digits);
  const normalized = !Number.isFinite(interval) || interval < 1
    ? 1
    : Math.min(Math.trunc(interval), 99);

  intervalDraft.value = String(normalized);
  emit("update-interval", normalized);
}

function normalizeTargetURL(rawURL) {
  const value = String(rawURL || "").trim();
  if (!value || !/^https?:\/\//i.test(value)) {
    return null;
  }

  try {
    const parsed = new URL(value);
    const host = parsed.hostname.replace(/^www\./, "").toLowerCase();
    if (!host || !host.includes(".")) {
      return null;
    }

    const port = parsed.port ? `:${parsed.port}` : "";
    parsed.protocol = parsed.protocol.toLowerCase();
    parsed.hostname = host;
    parsed.host = `${host}${port}`;

    if (!parsed.pathname || parsed.pathname === "/") {
      if (!parsed.search) {
        parsed.pathname = "/favicon.ico";
      } else {
        parsed.pathname = "/";
      }
    }

    return parsed.toString();
  } catch {
    return null;
  }
}

const isValidSiteURL = computed(() => {
  return Boolean(normalizeTargetURL(siteURL.value));
});

const isDuplicateSiteURL = computed(() => {
  const normalized = normalizeTargetURL(siteURL.value);
  if (!normalized) return false;
  return props.targetRows.some((target) => normalizeTargetURL(target.url) === normalized);
});

const siteURLError = computed(() => {
  if (!siteURLTouched.value || !siteURL.value.trim()) return "";
  if (!/^https?:\/\//i.test(siteURL.value.trim())) {
    return "URL must start with http:// or https://";
  }
  if (isDuplicateSiteURL.value) {
    return "This site URL already exists";
  }
  if (!isValidSiteURL.value) {
    return "Enter a valid site URL";
  }
  return "";
});

async function handleAddTarget() {
  siteURLTouched.value = true;
  if (!isValidSiteURL.value || isDuplicateSiteURL.value) return;
  const saved = await props.addTargetUrl(siteURL.value);
  if (!saved) return;
  siteURL.value = "";
  siteURLTouched.value = false;
  modalOpen.value = false;
}

const canDeleteTargets = computed(() => props.targetRows.length > 5);

function enterDeleteMode(targetID) {
  if (!canDeleteTargets.value) return;
  activeDeleteTargetID.value = activeDeleteTargetID.value === targetID ? "" : targetID;
}

async function handleRemoveTarget(targetID) {
  if (typeof props.removeTarget !== "function") {
    return;
  }
  const removed = await props.removeTarget(targetID);
  if (removed) {
    activeDeleteTargetID.value = "";
  }
}

function clearDeleteMode() {
  activeDeleteTargetID.value = "";
}

function handleGlobalPointerDown(event) {
  if (!activeDeleteTargetID.value) return;
  if (rootElement.value?.contains(event.target)) return;
  clearDeleteMode();
}

onMounted(() => {
  document.documentElement.addEventListener("pointerdown", handleGlobalPointerDown);
});

onBeforeUnmount(() => {
  document.documentElement.removeEventListener("pointerdown", handleGlobalPointerDown);
});
</script>

<template>
  <section ref="rootElement" class="flex flex-col gap-4" @click="clearDeleteMode">
    <div class="flex flex-wrap items-center justify-between gap-3 border-b pb-3">
      <div class="flex items-center gap-2">
        <Button
          size="sm"
          class="h-8 rounded-sm px-3 text-[13px] shadow-none active:translate-y-px cursor-pointer"
          :disabled="loading"
          @click="emit('benchmark')"
        >
          <Spinner v-if="loading" data-icon="inline-start" />
          <CirclePlay v-else data-icon="inline-start" />
          {{ loading ? "Benchmarking" : "Benchmark" }}
        </Button>
        <div class="flex h-8 items-center gap-2 rounded-sm border bg-background px-2.5">
          <Clock3 class="size-3.5 text-muted-foreground" />
          <input
            type="number"
            min="1"
            max="99"
            :value="intervalDraft"
            class="h-full w-9 border-0 bg-transparent px-0 text-right text-[13px] tabular-nums shadow-none outline-none ring-0 focus-visible:outline-none focus-visible:ring-0 focus-visible:ring-offset-0"
            @input="handleIntervalInput"
          />
          <span class="text-[12px] text-muted-foreground">m</span>
        </div>
        <Button
          size="sm"
          variant="secondary"
          class="h-8 rounded-sm px-2 shadow-none active:translate-y-px cursor-pointer"
          :disabled="saving"
          @click="modalOpen = true"
        >
          <Plus class="size-3.5" />
        </Button>
      </div>

      <div class="flex items-center gap-2">
        <Badge variant="secondary" class="h-6 rounded-sm px-2 text-[11px] font-medium">
          {{ availableTargets }}/{{ targetRows.length }} reachable
        </Badge>
        <Badge variant="secondary" class="h-6 rounded-sm px-2 text-[11px] font-medium">
          {{ averageLabel }}
        </Badge>
      </div>
    </div>

    <div class="overflow-hidden rounded-sm border bg-card">
      <div
        v-for="target in targetRows"
        :key="target.id"
        class="flex items-center gap-3 border-b px-3 py-2.5 last:border-b-0"
      >
        <button
          type="button"
          class="flex size-5 items-center justify-center cursor-pointer"
          :class="{ 'cursor-default': !canDeleteTargets }"
          @click.stop="enterDeleteMode(target.id)"
        >
          <Avatar
            v-if="activeDeleteTargetID !== target.id"
            class="size-5 rounded-sm"
            shape="square"
          >
            <AvatarImage :src="target.iconUrl" :alt="target.name" />
            <AvatarFallback>{{ fallbackLabel(target.name) }}</AvatarFallback>
          </Avatar>
          <svg
            v-else
            viewBox="0 0 1024 1024"
            xmlns="http://www.w3.org/2000/svg"
            class="size-5"
            @click.stop="handleRemoveTarget(target.id)"
          >
            <path
              d="M874.482336 149.501664c-199.356885-199.356885-525.623787-199.356885-724.980672 0s-199.356885 525.623787 0 724.980672 525.623787 199.356885 724.980672 0 199.356885-525.623787 0-724.980672z"
              fill="#FA453B"
            />
            <path
              d="M666.677583 739.18845L511.992 584.502867 357.306417 739.18845 284.79555 666.677583 439.481133 511.992 284.79555 357.306417 357.306417 284.79555 511.992 439.481133l154.685583-154.685583 72.510867 72.510867L584.502867 511.992l154.685583 154.685583z"
              fill="#FFFFFF"
            />
          </svg>
        </button>

        <div class="flex min-w-0 flex-1 items-center gap-3">
          <div class="w-36 truncate text-[13px] font-medium text-foreground">
            {{ target.name }}
          </div>

          <div class="min-w-0 flex-1 truncate text-[12px] text-muted-foreground">
            {{ target.url }}
          </div>
        </div>

        <Badge class="h-6 rounded-sm px-2 text-[11px] tabular-nums" :variant="target.statusVariant">
          {{ target.statusText }}
        </Badge>
      </div>
    </div>

    <Modal v-model:open="modalOpen" backdrop="blur" size="sm">
      <div class="flex flex-col gap-1">
        <DialogTitle class="text-[15px] font-semibold tracking-tight">Add Checkpoint</DialogTitle>
        <p class="text-[12px] text-muted-foreground">Enter a website URL to monitor its connectivity.</p>
      </div>

      <div class="mt-2 flex flex-col gap-2">
        <Label for="custom-site-url" class="text-[12px] font-medium">Site URL</Label>
        <Input
          id="custom-site-url"
          v-model="siteURL"
          placeholder="https://github.com"
          class="h-9 rounded-sm text-[13px] shadow-none"
          @input="siteURLTouched = true"
          @keyup.enter="isValidSiteURL && !isDuplicateSiteURL && handleAddTarget()"
        />
        <p v-if="siteURLError" class="text-[11px] text-destructive">
          {{ siteURLError }}
        </p>
      </div>

      <div class="mt-4 flex justify-end gap-2">
        <Button
          size="sm"
          variant="ghost"
          class="h-8 rounded-sm px-3 text-[13px]"
          @click="modalOpen = false"
        >
          Cancel
        </Button>
        <Button
          size="sm"
          class="h-8 rounded-sm px-3 text-[13px] shadow-none"
          :disabled="saving || !isValidSiteURL || isDuplicateSiteURL"
          @click="handleAddTarget"
        >
          {{ saving ? "Saving" : "Add Checkpoint" }}
        </Button>
      </div>
    </Modal>
  </section>
</template>

<style scoped>
</style>
