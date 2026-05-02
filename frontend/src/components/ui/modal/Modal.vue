<script setup>
import { DialogContent, DialogOverlay, DialogPortal, DialogRoot } from "reka-ui";
import { cn } from "@/lib/utils";

const props = defineProps({
  open: { type: Boolean, default: false },
  size: { type: String, default: "sm" }, // xs, sm, md, lg
  backdrop: { type: String, default: "blur" }, // blur, opaque, transparent
});

const emit = defineEmits(["update:open"]);

const sizeClasses = {
  xs: "max-w-xs",
  sm: "max-w-sm",
  md: "max-w-md",
  lg: "max-w-lg",
};

const backdropClasses = {
  blur: "bg-black/20 backdrop-blur-md",
  opaque: "bg-black/50",
  transparent: "bg-transparent",
};
</script>

<template>
  <DialogRoot :open="open" @update:open="(val) => emit('update:open', val)">
    <DialogPortal>
      <DialogOverlay
        :class="cn(
          'fixed inset-0 z-50 flex items-center justify-center p-4 transition-all duration-200',
          backdropClasses[backdrop]
        )"
      />
      <DialogContent
        :class="cn(
          'fixed left-[50%] top-[50%] z-50 flex w-full translate-x-[-50%] translate-y-[-50%] flex-col gap-4 rounded-2xl bg-card p-6 shadow-xl outline-none transition-all duration-200 sm:rounded-3xl',
          sizeClasses[size]
        )"
      >
        <slot />
      </DialogContent>
    </DialogPortal>
  </DialogRoot>
</template>
