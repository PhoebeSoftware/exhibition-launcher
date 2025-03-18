<template>
  <div class="sidebar">
    <div class="sidebar-content">
      <router-link to="/explore" custom v-slot="{ navigate, isActive }">
        <button @click="navigate" class="nav-button">
          <div v-if="isActive" class="active-indicator-sidebar"></div>
          Explore
        </button>
      </router-link>
      <router-link to="/library" custom v-slot="{ navigate, isActive }">
        <button @click="navigate" class="nav-button">
          <div v-if="isActive" class="active-indicator-sidebar"></div>
          Library
        </button>
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
  padding: 0 20px;
  height: calc(100% - 40px);
  background-color: rgb(22, 22, 22);
  position: fixed;
  bottom: 0;
}

.active-indicator-sidebar {
  width: 4px;
  height: 0;
  background-color: var(--accent-color);
  box-shadow: 0px 0px 20px var(--accent-color);
  border-radius: 50px;
  position: absolute;
  left: 0;
  margin-left: -10px;
  top: 50%;
  transform: translateY(-50%);
  animation: growHeight 0.3s ease forwards;
}

@keyframes growHeight {
  from {
    height: 0;
  }
  to {
    height: 35px;
  }
}

.sidebar-content {
  padding-top: 30px;
  width: 100%;
  height: auto;
  display: flex;
  flex-direction: column;
}

.nav-button {
  position: relative;
  color: white;
  background: none;
  text-decoration: none;
  border: none;
  padding: 15px 20px;
  text-align: left;
  border-radius: 10px;
  transition: background-color 0.3s ease;
  cursor: pointer;
}

.nav-button:hover {
  background-color: rgba(255, 255, 255, 0.025);
}
</style>