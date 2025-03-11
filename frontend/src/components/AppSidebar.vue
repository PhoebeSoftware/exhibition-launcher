<template>
  <div class="sidebar">
    <div class="sidebar-content">
      <router-link to="/explore" custom v-slot="{ navigate }">
        <button @click="navigate">Explore</button>
      </router-link>
      <router-link to="/library" custom v-slot="{ navigate }">
        <button @click="navigate">Library</button>
      </router-link>
    </div>
  </div>
</template>

<script>
export default {
  name: 'AppSidebar',
}
</script>

<style scoped>
.sidebar {
  width: 200px;
  height: calc(100% - 40px);
  background-color: rgb(68, 48, 48);
  position: fixed;
  bottom: 0;
}

.sidebar-content {
  padding-top: 30px;
  width: 100%;
  height: auto;
  display: flex;
  flex-direction: column;
}

.sidebar-content button {
  color: white;
  background: none;
  text-decoration: none;
  border: none;
  padding: 10px;
  text-align: left;
  transition: background-color 0.3s ease;
  cursor: pointer;
}

.sidebar-content button:hover {
  background-color: rgba(255, 255, 255, 0.164);
}
</style>