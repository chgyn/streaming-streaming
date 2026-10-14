<script setup>
import Navbar from './components/Navbar.vue'
import PlayerModal from './components/PlayerModal.vue'
import Home from './pages/Home.vue'
</script>

<template>
  <div class="min-h-screen flex flex-col bg-slate-950 text-slate-100 selection:bg-sky-500 selection:text-slate-950">
    <Navbar />
    <div class="flex-1">
      <Home />
    </div>
    <PlayerModal />
  </div>
</template>
