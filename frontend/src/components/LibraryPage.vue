<template>
    <div class="page">
        <!-- Library Page -->
        <div v-if="currentPage === 'library'">
            <div class="Library-settings-container">
                <div class="Library-settings-wrapper">
                    <span>Sort by:</span>
                    <button id="sort-library-by">Time played <i class="fa-solid fa-chevron-down"></i></button>
                </div>
                <div class="Library-favorites-container">
                    <button id="all-button" class="active">
                        All
                        <div class="active-indicator-horizontal"></div>
                    </button>
                    <button id="favorites-button">
                        Favorites
                        <div class="active-indicator-horizontal"></div>
                    </button>
                </div>
                <div class="add-game-to-library-wrapper">
                    <button id="refresh-page"><i class="fa-solid fa-arrows-rotate"></i></button>
                    <button class="game-add-button" id="Import-from-PC-TT" @click="addGame">
                        Import from PC<i class="fa-solid fa-desktop"></i>
                    </button>
                </div>
            </div>

            <div class="game-library-container">
                <div class="game-library-game-box" v-for="game in games" :key="game.igdb_id"
                     @click="openGameStore(game.igdb_id)"
                     :style="{ backgroundImage: `url(${game.MainCover})`}">
                    <div class="game-box-info">
                        <div class="text-container">
                            <h1>{{ game.name }}</h1>
                            <p>{{ game.executable }}</p>
                        </div>
                        <button><i class="fa-solid fa-ellipsis"></i></button>
                    </div>
                </div>
            </div>
        </div>

        <!-- Simple Game Store Page -->
        <div v-if="currentPage === 'store'" class="game-store-page">
            <div class="game-store-header">
                <button class="back-button" @click="returnToLibrary">
                    <i class="fa-solid fa-arrow-left"></i> Back to Library
                </button>
            </div>

            <div class="game-store-content">
                <h1>Game Store Page</h1>
                <p>This is a placeholder for game #{{ selectedGameId }}</p>
                <p>Store page content will be implemented later</p>
            </div>
        </div>
    </div>
</template>

<script>
import {Library} from '../../bindings/derpy-launcher072/library';
import router from "@/router.js";

export default {
    name: 'LibraryPage',

    data() {
        return {
            games: [],
            currentPage: 'library',
            selectedGameId: null
        };
    },
    methods: {
        addGame() {
            Library.AddToLibrary(11544).catch(console.warn); // ELDEN RING ID
            this.router.go()
        },

        openGameStore(gameId) {
            this.selectedGameId = gameId;
            this.currentPage = 'store';
            console.log(`Opening store page for game ${gameId}`);
        },

        returnToLibrary() {
            this.currentPage = 'library';
        }
    },
    async mounted() {
        Library.GetAllGames().then((games) => {
            Object.values(games).forEach((game) => {
                console.log(game);
                this.games.push(game);
            });
        });
    }
};
</script>

<style scoped>
.page {
    padding: 20px;
    height: auto;
}

.game-library-container {
    width: 100%;
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(170px, 2fr));
    gap: 20px;
    padding-top: 30px;
    margin-bottom: 50px;
}

.add-game-to-library-wrapper {
    display: flex;
    justify-content: center;
    align-items: center;
}

#refresh-page i {
    color: var(--secondary-text-color);
    transition: color 0.3s ease;
}

#refresh-page:hover i {
    color: var(--text-color);
}

#refresh-page {
    background: none;
    border: none;
}

.add-game-to-library-wrapper button {
    outline: none;
    border: none;
    padding: 10px 20px;
    color: var(--secondary-text-color);
    border-radius: 15px;
    background-color: var(--hover-background-color);
    border: 1px solid var(--outline);
    cursor: pointer;
    transition: all 0.2s ease;
    display: flex;
    gap: 10px;
}

.add-game-to-library-wrapper button:hover {
    color: var(--text-color);
}

.Library-settings-wrapper {
    display: flex;
    justify-content: left;
    align-items: center;
    width: auto;
    gap: 15px;
    font-size: 16px;
    color: var(--secondary-text-color);
}

.Library-settings-wrapper button {
    border: none;
    background: none;
    cursor: pointer;
    color: var(--secondary-text-color);
    padding: 10px 20px;
    border-radius: 15px;
    background-color: var(--hover-background-color);
    border: 1px solid var(--outline);
}

.Library-favorites-container {
    width: auto;
    height: auto;
    display: flex;
    gap: 20px;
    position: relative;
}

.Library-favorites-container button {
    background: none;
    border: none;
    color: var(--secondary-text-color);
    font-size: 15px;
    cursor: pointer;
    padding: 5px 10px;
    position: relative;
    transition: color 0.3s ease, background-color 0.3s ease;
    border-radius: 15px;
}

.Library-favorites-container button:hover {
    background-color: rgb(29, 29, 29);
}

.Library-favorites-container button.active {
    color: var(--accent-color);
    background-color: rgba(0, 0, 0, 0) !important;
}

.Library-favorites-container button.active .active-indicator-horizontal {
    opacity: 1;
    animation: growWidth 0.3s ease forwards;
}

.Library-settings-container {
    width: 100%;
    height: 90px;
    display: flex;
    justify-content: space-between;
    align-items: center;
    position: sticky;
    top: 50px;
    padding: 10px;
    backdrop-filter: blur(15px);
    background-color: rgba(25, 25, 25, 0.4);
    z-index: 1;
}

.game-library-game-box {
    width: 100%;
    height: 250px;
    border-radius: 15px;
    background-color: var(--hover-background-color);
    display: block;
    flex-direction: column;
    justify-content: space-between;
    overflow: hidden;
    position: relative;
    cursor: pointer;
    transition: transform 0.3s ease;
    background-position: center; /* Center the image */
    background-repeat: no-repeat; /* Do not repeat the image */
    background-size: cover; /* Resize the background image to cover the entire container */
}

.game-library-game-box[style*="display: none"] {
    display: none !important;
}

.game-box-info {
    width: 100%;
    background-color: var(--game-box-info-background-color);
    color: var(--text-color);
    font-size: 13px;
    position: absolute;
    bottom: -100px;
    padding: 10px 20px;
    transition: all 0.3s ease;
    display: flex;
    align-items: center;
    justify-content: space-between;
}

.game-box-info p {
    color: var(--secondary-text-color);
    margin: 0;
}

.game-box-info h1 {
    font-size: 15px;
    margin: 0;
}

.game-library-game-box:hover .game-box-info {
    bottom: 0;
}

.game-box-info .text-container {
    display: flex;
    flex-direction: column;
    gap: 4px;
}

.game-box-info h1 {
    font-size: 12px;
}

.game-box-info i {
    transition: color 0.3s ease, transform 0.3s ease;
    color: var(--secondary-text-color);
    font-size: 20px;
}

.game-box-info button {
    right: 60px;
    border: none;
    background: none;
    cursor: pointer;
}

/* New store page styles */
.game-store-page {
    width: 100%;
    height: 100%;
}

.game-store-header {
    margin-bottom: 20px;
}

.back-button {
    background: none;
    border: none;
    color: var(--secondary-text-color);
    cursor: pointer;
    padding: 10px;
    border-radius: 15px;
    background-color: var(--hover-background-color);
    border: 1px solid var(--outline);
    transition: all 0.2s ease;
}

.back-button:hover {
    color: var(--text-color);
}

.game-store-content {
    padding: 20px;
    background-color: var(--hover-background-color);
    border-radius: 15px;
}
</style>