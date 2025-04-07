<template>
    <div class="sidebar">
        <div class="sidebar-content">
            <router-link to="/explore" custom v-slot="{ navigate, isActive }">
                <button @click="navigate" class="nav-button" :class="{ 'nav-button-active': isActive }">
                    <div v-if="isActive" class="active-indicator-sidebar"></div>
                    <i class="fas fa-compass"></i>
                    <span>Explore</span>
                </button>
            </router-link>
            <router-link to="/library" custom v-slot="{ navigate, isActive }">
                <button @click="navigate" class="nav-button" :class="{ 'nav-button-active': isActive }">
                    <div v-if="isActive" class="active-indicator-sidebar"></div>
                    <i class="fas fa-th-large"></i>
                    <span>Library</span>
                </button>
            </router-link>
            
            <div class="section-title">YOUR GAMES</div>
            
            <div class="sidebar-games-container">
                <div class="sidebar-game-item">
                    <div class="sidebar-game-icon sidebar-game"></div>
                    <span class="sidebar-game-title">Dark Souls 3</span>
                </div>
                
                <div class="sidebar-game-item">
                    <div class="sidebar-game-icon sidebar-game"></div>
                    <span class="sidebar-game-title">Sekiro</span>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
export default {
    name: 'AppSidebar'
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

@keyframes growHeight {
    from {
        height: 0;
    }
    to {
        height: 35px;
    }
}

.active-indicator-sidebar {
    width: 4px;
    height: 0;
    background-color: var(--accent-color);
    box-shadow: 0px 0px 20px var(--accent-color);
    border-radius: 0 4px 4px 0;
    position: absolute;
    left: -20px;
    top: 50%;
    transform: translateY(-50%);
    animation: growHeight 0.3s ease forwards;
}

.sidebar-content {
    padding-top: 30px;
    width: 100%;
    height: auto;
    display: flex;
    flex-direction: column;
    gap: 4px;
}

.nav-button {
    position: relative;
    display: flex;
    align-items: center;
    color: rgba(255, 255, 255, 0.4);
    background: none;
    border: none;
    padding: 10px 16px;
    text-align: left;
    border-radius: 6px;
    transition: all 0.2s ease;
    cursor: pointer;
    width: 100%;
    font-size: 14px;
}

.nav-button:hover {
    background-color: rgba(255, 255, 255, 0.05);
    color: rgba(255, 255, 255, 0.8);
}

.nav-button-active {
    background-color: rgba(255, 255, 255, 0.08) !important;
    color: white !important;
    font-weight: 500;
}

.nav-button i {
    margin-right: 12px;
    font-size: 16px;
    width: 20px;
    text-align: center;
    color: rgba(255, 255, 255, 0.4);
    transition: color 0.2s ease;
}

.nav-button:hover i {
    color: rgba(255, 255, 255, 0.8);
}

.nav-button-active i {
    color: rgba(255, 255, 255, 0.9);
}

.section-title {
    font-size: 12px;
    color: #6b7280;
    margin-top: 20px;
    padding: 8px 16px;
    font-weight: 500;
    letter-spacing: 0.5px;
}

.sidebar-games-container {
    display: flex;
    flex-direction: column;
    width: 100%;
    gap: 4px;
    margin-top: 4px;
}

.sidebar-game-item {
    display: flex;
    align-items: center;
    padding: 8px 12px;
    border-radius: 6px;
    cursor: pointer;
    transition: all 0.2s ease;
    color: rgba(255, 255, 255, 0.5);
}

.sidebar-game-item:hover {
    background-color: rgba(255, 255, 255, 0.05);
    color: rgba(255, 255, 255, 0.8);
}

.sidebar-game-icon {
    width: 32px;
    height: 32px;
    border-radius: 6px;
    background-position: center;
    background-size: cover;
    background-repeat: no-repeat;
    margin-right: 12px;
    flex-shrink: 0;
}

.sidebar-game {
    background-image: url("https://images.igdb.com/igdb/image/upload/t_cover_big/co698n.jpg");
}

.sidebar-game-title {
    font-size: 13px;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    transition: inherit;
}
</style>
