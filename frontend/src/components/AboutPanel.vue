<script setup>
const links = [
  { label: 'Wails v3', href: 'https://v3alpha.wails.io/' },
  { label: 'IPCheck.ing', href: 'https://ipcheck.ing/' },
]
</script>

<template>
  <section class="panel about">
    <header class="about__header">
      <img class="about__logo" src="/logo.svg" alt="" aria-hidden="true">
      <div>
        <h1 class="panel__title">SiteCheck</h1>
        <p class="panel__subtitle">Menubar-only network connectivity monitor for macOS.</p>
      </div>
    </header>

    <div class="about__body">
      <dl class="about__facts">
        <div>
          <dt>Version</dt>
          <dd>0.1.0</dd>
        </div>
        <div>
          <dt>Runtime</dt>
          <dd>Wails v3 alpha</dd>
        </div>
        <div>
          <dt>Probe</dt>
          <dd>GET favicon endpoint, no-store cache policy</dd>
        </div>
      </dl>

      <div class="about__links">
        <a
          v-for="link in links"
          :key="link.href"
          class="button button--outline"
          href="#"
          :data-wml-openURL="link.href"
        >
          {{ link.label }}
        </a>
      </div>
    </div>
  </section>
</template>
