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
                     @click="openGamePage(game)"
                     :style="{
                         backgroundImage: `url(${getCoverUrlFromGame(game)})`
                     }">
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
                <div>
                    <div class="game-store-banner"></div>
                    <p>This is a placeholder for game #{{ selectedGame.igdb_id }}</p>
                    <p>Store page content will be implemented later</p>
                </div>
            </div>
        </div>

        <!-- Game Page -->
        <div v-if="currentPage === 'game'" class="game-page">
            <div class="game-page-content">
                <button class="back-button" @click="returnToLibrary">
                    <i class="fa-solid fa-arrow-left"></i> Back to Library
                </button>
                <!-- Add more game details here as needed -->
            </div>
            <div class="game-page-image-container"
                 :style="{ backgroundImage: `url(${getBackgroundImage()})` }">
                <div class="game-user-stats">
                    <div class="game-user-stats-left">
                        <button @click="launchGame(selectedGame)">
                            <i class="fa-solid fa-play"></i>PLAY
                        </button>
                        <div class="last-played-wrapper">
                            <h1>Last played</h1>
                            <p>{{ selectedGame.last_played || 'Never' }}</p>
                        </div>
                        <div class="play-time-wrapper">
                            <h1>Play time</h1>
                            <p>{{ selectedGame.playTime || '0 hours' }}</p>
                        </div>
                    </div>

                    <div class="game-options-wrapper">
                        <button @click="openGameSettings">
                            <i class="fa-solid fa-gear"></i>
                        </button>

                        <button @click="toggleFavorite">
                            <i :class="selectedGame.isFavorite ? 'fa-solid fa-star' : 'fa-regular fa-star'"></i>
                        </button>
                    </div>
                    <div>
                        <p>{{ selectedGame.description }}</p>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import {Library} from "../../bindings/exhibition-launcher/library/index.js";

export default {
    name: 'LibraryPage',

    data() {
        return {
            games: [],
            currentPage: 'library',
            selectedGame: null,
        };
    },
    methods: {
        getBackgroundImage() {
            let url;
            if (this.selectedGame.artwork_url_list != null && this.selectedGame.artwork_url_list.length > 0) {
                url = this.selectedGame.artwork_url_list[0];
            } else if (this.selectedGame.screenshot_url_list != null && this.selectedGame.screenshot_url_list.length > 0) {
                url = this.selectedGame.screenshot_url_list[0];
            } else {
                url = this.selectedGame.cover_url;
            }
            return new URL(url, import.meta.url).href
        },
        async addGame() {
            let newGame = await Library.AddToLibrary(119277, true).catch((err) => {
                console.warn(err)
            });
            this.games.push(newGame)
        },

        getCoverUrlFromGame(game) {
            return new URL(game.cover_url, import.meta.url).href
        },

        openGameStore(game) {
            this.selectedGame = game;
            this.currentPage = 'store';
        },

        openGamePage(game) {
            this.selectedGame = game;
            this.currentPage = 'game';
        },

        launchGame(game) {
            // Implement game launch logic
            Library.StartApp(game.igdb_id).catch((err) => {
                console.log(err)
            })
        },

        openGameSettings() {
            // Implement game settings logic
            console.log(`Opening settings for game ${this.selectedGame.igdb_id}`);
        },

        toggleFavorite() {
            // Toggle favorite status
            if (this.selectedGame) {
                this.selectedGame.isFavorite = !this.selectedGame.isFavorite;
                console.log(`Favorite status for game ${this.selectedGame.igdb_id}: ${this.selectedGame.isFavorite}`);
            }
        },

        returnToLibrary() {
            this.currentPage = 'library';
        }
    },
    async mounted() {
        const amountOfGames = await Library.GetAmountOfGames()

        console.log(amountOfGames)
        const portion = 100;

        for (let i = 0; i < amountOfGames; i += portion) {
            console.log(i)
            let games = await Library.GetRangeGame(portion, i)
            for (let j = 0; j < games.length; j++) {
                let game = games[j]
                console.log(game.name + " : " + game.igdb_id + " : " + j);
                this.games.push(game);
            }
        }

        // Library.GetRangeGame(100, 0).then()
        //
        //
        // Library.GetAllGames().then((games) => {
        //     Object.values(games).forEach((game) => {
        //         console.log(game);
        //         this.games.push(game);
        //     });
        // });
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
    background-position: center;
    background-repeat: no-repeat;
    background-size: cover;
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

/* Game Store Page Styles */
.game-store-page {
    width: 100%;
    height: 100%;
}

.game-store-header {
    margin-bottom: 20px;
}

.game-store-content {
    padding: 20px;
    background-color: var(--hover-background-color);
    border-radius: 15px;
}

.game-store-banner {
    height: 900px;
    width: 100%;
    background-position: center;
    background-repeat: no-repeat;
    background-size: cover;
}

/* Game Page Styles */
.game-page {
    position: relative;
    width: 100%;
    background-color: var(--background-color);
}

.game-page-content {
    padding: 20px;
}

.game-page-image-container {
    width: 100%;
    height: 500px;
    background-color: rgb(35, 35, 35);
    background-position: center;
    background-repeat: no-repeat;
    border-radius: 15px;
    display: flex;
    align-items: end;
    position: relative;
}

.game-user-stats-container {
    position: sticky;
    top: 0;
    left: 0;
    width: 100%;
    z-index: 10;
}

.game-user-stats {
    width: 100%;
    position: sticky;
    top: 0;
    left: 0;
    backdrop-filter: blur(10px);
    background: linear-gradient(
        to bottom,
        rgba(25, 25, 25, 0.6) 0%,
        rgba(25, 25, 25, 0.8) 50%,
        rgba(25, 25, 25, 1) 100%
    );
    display: flex;
    flex-wrap: wrap;
    align-items: flex-start;
    justify-content: space-between;
    padding: 10px 20px;
    gap: 10px;
}

.game-user-stats p {
    color: var(--secondary-text-color);
}

.last-played-wrapper {
    display: flex;
    flex-direction: column;
    color: var(--text-color);
    font-size: 12px;
}

.play-time-wrapper {
    font-size: 12px;
    display: flex;
    flex-direction: column;
    color: var(--text-color);
}

.game-user-stats-left {
    display: flex;
    gap: 20px;
    flex-direction: row;
    align-items: center;
}

.game-user-stats-left button {
    cursor: pointer;
    padding: 15px 40px;
    border: 0;
    background-color: var(--accent-color);
    color: var(--text-color);
    display: flex;
    align-items: center;
    gap: 10px;
    border-radius: 10px;
}

.game-options-wrapper {
    display: flex;
    gap: 10px;
}

.game-options-wrapper button {
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

.game-options-wrapper button:hover {
    color: var(--text-color);
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

.placeholder {
    height: 5000px;
    width: 100%;
    background-color: orange;
}
</style>