<template>
    <footer class="app-footer">
      <h1>{{ placeholderTitle }}</h1>
      <div class="footer-content">
        <button 
          v-if="!isDownloading" 
          id="placeholder-text" 
          @click="toggleDownload"
        >
          <i class="fa-solid fa-download"></i>
          No downloads in progress
        </button>
        
        <div 
          v-else 
          class="download-item" 
          id="download-item"
        >
          <img 
            id="download-image" 
            :src="downloadImage" 
            alt="Download Image"
          >
          <div class="download-content">
            <div class="info-container">
              <p class="game-name">{{ gameName }}</p>
              <p class="progress-percentage">{{ progressPercentage }}%</p>
            </div>
            <div class="progress-bar">
              <div 
                class="progress-bar-fill" 
                :style="{ width: `${progressPercentage}%` }"
              ></div>
            </div>
            <div class="download-speed">
              <p><span id="download-speed">{{ downloadSpeed }} MB/s</span></p>
            </div>
          </div>
        </div>
      </div>
      <h1>{{ version }}</h1>
    </footer>
  </template>
  
  <script>
  export default {
    name: 'AppFooter',
    data() {
      return {
        placeholderTitle: 'Placeholder',
        version: 'v0.0.1',
        isDownloading: false,
        downloadImage: 'https://steamuserimages-a.akamaihd.net/ugc/687094810512264399/04BA8A55B390D1ED0389E561E95775BCF33A9857/?imw=5000&imh=5000&ima=fit&impolicy=Letterbox&imcolor=%23000000&letterbox=false',
        gameName: 'Rust',
        progressPercentage: 69,
        downloadSpeed: 0
      }
    },
    methods: {
      toggleDownload() {
        this.isDownloading = !this.isDownloading;
      }
    }
  }
  </script>
  
  <style scoped>
  .download-item {
    display: flex;
    flex-direction: row;
    align-items: center;
    gap: 10px;
    cursor: pointer;
    overflow: hidden;
    animation: ExpandDownload 0.3s ease forwards;
  }
  
  @keyframes ExpandDownload {
    from { min-width: 10%; }
    to { min-width: 40%; }
  }
  
  .download-item img {
    width: 35px;
    height: 35px;
    border-radius: 10px;
    object-fit: cover;
  }
  
  .download-content {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 5px;
  }
  
  .info-container {
    width: 100%;
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
  
  .game-name {
    color: var(--text-color);
    font-size: 12px;
    font-weight: 400;
    margin: 0;
  }
  
  .progress-percentage {
    color: var(--text-color);
    font-size: 12px;
    margin: 0;
  }
  
  .download-speed {
    color: var(--secondary-text-color);
    font-size: 12px;
    margin: 0;
  }
  
  .progress-bar {
    width: 100%;
    height: 2px;
    background-color: var(--progress-bar-background-color);
    position: relative;
    border-radius: 3px;
  }
  
  .progress-bar-fill {
    height: 100%;
    background-color: var(--secondary-text-color);
    box-shadow: 0px 0px 20px var(--secondary-text-color);
    position: absolute;
    border-radius: 3px;
    transition: width 0.3s linear, color 0.3s linear;
    top: 0;
    left: 0;
  }
  
  .app-footer {
    width: 100%;
    height: 60px;
    background-color: rgb(22, 22, 22);
    position: fixed;
    bottom: 0;
    display: flex;
    justify-content: space-between;
    align-items: center;
    color: var(--secondary-text-color);
    padding: 0 20px;
    box-sizing: border-box;
  }
  
  .app-footer h1 {
    font-size: 14px;
  }
  
  .footer-content {
    flex: 1;
    display: flex;
    justify-content: center;
    align-items: center;
  }
  
  #placeholder-text {
    display: flex;
    font-size: 16px;
    color: var(--secondary-text-color);
    gap: 10px;
    background: none;
    border: none;
    cursor: pointer;
  }
  
  #download-item {
    display: flex;
    align-items: center;
    gap: 10px;
  }
  </style>