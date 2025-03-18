<template>
    <div class="topbar">
        <div class="topbar-left">
            <span class="app-name">Sigma Launcher</span>
        </div>
        <div class="topbar-right">
            <button class="topbar-button" id="minimize-btn" @click="minimize">
                <i class="fa-solid fa-minus"></i>
            </button>

            <button class="topbar-button" id="maximize-btn" @click="maximize">
                <i class="fa-regular fa-square"></i>
            </button>

            <button class="topbar-button" id="close-btn" @click="close">
                <i class="fa-solid fa-xmark"></i>
            </button>
        </div>
    </div>
</template>

<script>
import { Close, Maximize, Minimize } from '../../bindings/derpy-launcher072/windowservice';

export default {
    name: 'AppTopbar',
    methods: {
        minimize() {
            console.log('Minimize button clicked');
            Minimize();
        },
        maximize() {
            console.log('Maximize button clicked');
            Maximize();
        },
        close() {
            console.log('Close button clicked');
            Close();
        }
    }
}
</script>

<style scoped>
.topbar {
    --wails-draggable: drag;
    z-index: 10;
    position: fixed;
    color: white;
    top: 0;
    left: 0;
    width: 100%;
    background-color: rgb(22, 22, 22);
    padding: 10px 20px;
    display: flex;
    justify-content: space-between;
    align-items: center;
    box-sizing: border-box;
}

.topbar-right {
    display: flex;
    flex-direction: row;
    gap: 10px;
}

.topbar-button {
    background: none;
    border: none;
    color: white;
    font-size: 16px;
    cursor: pointer;
    padding: 5px 10px;
    border-radius: 4px;
    transition: background-color 0.3s ease, color 0.3s ease;
    display: flex;
    justify-content: center;
    align-items: center;
    width: 30px;
    height: 30px;
}

.topbar-button:hover {
    background-color: rgb(29, 29, 29);
    color: var(--text-color);
}

.topbar-button i {
    font-size: 12px;
    margin: 0;
}
</style>