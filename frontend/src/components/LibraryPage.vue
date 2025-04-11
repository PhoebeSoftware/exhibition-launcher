<template>
    <div class="page">
        <!-- Library Page -->
        <div v-if="currentPage === 'library'">
            <div class="Library-settings-container">
                <div class="Library-settings-wrapper">
                    <span>Sort by:</span>
                    <button id="sort-library-by">
                        Time played <i class="fa-solid fa-chevron-down"></i>
                    </button>
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
                    <button id="refresh-page">
                        <i class="fa-solid fa-arrows-rotate"></i>
                    </button>
                    <button
                        class="game-add-button"
                        id="Import-from-PC-TT"
                        @click="addGame"
                    >
                        Import from PC<i class="fa-solid fa-desktop"></i>
                    </button>
                </div>
            </div>

            <div class="game-library-container">
                <div
                    class="game-library-game-box"
                    v-for="game in games"
                    :key="game.igdb_id"
                    @click="openGamePage(game)"
                    :style="{
                        backgroundImage: `url(${coverUrls[game.igdb_id]})`,
                    }"
                >
                    <div class="game-box-info">
                        <div class="text-container">
                            <h1>{{ game.name }}</h1>
                        </div>
                        <button><i class="fa-solid fa-ellipsis"></i></button>
                    </div>
                </div>
            </div>
        </div>

        <div v-if="currentPage === 'game'" class="game-page-new">
            <button class="back-button" @click="returnToLibrary">
                <i class="fa-solid fa-arrow-left"></i>
            </button>

            <!-- Image Carousel Section -->
            <div class="carousel-container">
                <div
                    class="carousel-slides"
                    ref="carouselSlides"
                    :style="{
                        transform: `translateX(-${currentSlide * 100}%)`,
                    }"
                >
                    <div
                        class="carousel-slide"
                        v-for="(image, index) in getGameImages()"
                        :key="index"
                        :style="{ backgroundImage: `url(${image})` }"
                    >
                        <div class="banner-overlay"></div>
                    </div>
                </div>

                <div class="carousel-controls">
                    <button class="carousel-btn prev" @click="prevSlide">
                        <i class="fa-solid fa-chevron-left"></i>
                    </button>
                    <div class="carousel-indicators">
                        <span
                            v-for="(_, index) in getGameImages()"
                            :key="index"
                            :class="{ active: index === currentSlide }"
                            @click="goToSlide(index)"
                        >
                        </span>
                    </div>
                    <button class="carousel-btn next" @click="nextSlide">
                        <i class="fa-solid fa-chevron-right"></i>
                    </button>
                </div>

                <div class="game-info-overlay">
                    <div class="game-branding">
                        <h1 class="game-title">{{ selectedGame.name }}</h1>
                    </div>

                    <div class="game-stats-container">
                        <div class="stat-box">
                            <div class="stat-value">
                                {{ selectedGame.playTime || "16hrs" }}
                                <i class="fa-regular fa-clock"></i>
                            </div>
                            <div class="stat-label">PLAY TIME</div>
                        </div>

                        <div class="stat-box">
                            <div class="stat-value">
                                {{ selectedGame.achievementPercent || "27%" }}
                                <i class="fa-solid fa-trophy"></i>
                            </div>
                            <div class="stat-label">ACHIEVEMENTS</div>
                        </div>
                    </div>

                    <div class="game-actions">
                        <button class="add-library-btn" @click="toggleFavorite">
                            {{
                                selectedGame.isFavorite
                                    ? "Remove from Favorites"
                                    : "Add to Favorites"
                            }}
                            <i
                                :class="
                                    selectedGame.isFavorite
                                        ? 'fa-solid fa-heart'
                                        : 'fa-regular fa-heart'
                                "
                            ></i>
                        </button>
                        <button
                            class="download-btn"
                            @click="launchGame(selectedGame)"
                        >
                            Play Game <i class="fa-solid fa-play"></i>
                        </button>
                    </div>
                </div>
            </div>

            <div class="game-details-section">
                <div class="game-details-container">
                    <div class="details-column">
                        <div class="details-section">
                            <h3>About</h3>
                            <p>
                                {{
                                    selectedGame.fullDescription ||
                                    selectedGame.description ||
                                    "A groundbreaking game experience."
                                }}
                            </p>
                        </div>

                        <div class="details-section">
                            <h3>Genre</h3>
                            <div class="tags-container">
                                <span
                                    class="tag"
                                    v-for="(
                                        genre, index
                                    ) in selectedGame.genres || [
                                        'RPG',
                                        'Fantasy',
                                        'Adventure',
                                    ]"
                                    :key="index"
                                >
                                    {{ genre }}
                                </span>
                            </div>
                        </div>

                        <div class="details-section">
                            <h3>Features</h3>
                            <div class="tags-container">
                                <span
                                    class="tag"
                                    v-for="(
                                        feature, index
                                    ) in selectedGame.features || [
                                        'Single-player',
                                        'Controller Support',
                                        'Cloud Saves',
                                    ]"
                                    :key="index"
                                >
                                    {{ feature }}
                                </span>
                            </div>
                        </div>
                    </div>

                    <div class="details-column">
                        <div class="details-section">
                            <h3>System Requirements</h3>
                            <div class="system-reqs">
                                <div class="req-section">
                                    <h4>Minimum</h4>
                                    <ul>
                                        <li>
                                            <strong>OS:</strong>
                                            {{
                                                selectedGame.minOS ||
                                                "Windows 10 (64-bit)"
                                            }}
                                        </li>
                                        <li>
                                            <strong>CPU:</strong>
                                            {{
                                                selectedGame.minCPU ||
                                                "Intel Core i5-2500K | AMD FX-8320"
                                            }}
                                        </li>
                                        <li>
                                            <strong>RAM:</strong>
                                            {{ selectedGame.minRAM || "8 GB" }}
                                        </li>
                                        <li>
                                            <strong>GPU:</strong>
                                            {{
                                                selectedGame.minGPU ||
                                                "NVIDIA GeForce GTX 760 | AMD Radeon HD 7950"
                                            }}
                                        </li>
                                        <li>
                                            <strong>Storage:</strong>
                                            {{
                                                selectedGame.minStorage ||
                                                "60 GB available space"
                                            }}
                                        </li>
                                    </ul>
                                </div>

                                <div class="req-section">
                                    <h4>Recommended</h4>
                                    <ul>
                                        <li>
                                            <strong>OS:</strong>
                                            {{
                                                selectedGame.recOS ||
                                                "Windows 10/11 (64-bit)"
                                            }}
                                        </li>
                                        <li>
                                            <strong>CPU:</strong>
                                            {{
                                                selectedGame.recCPU ||
                                                "Intel Core i7-4790 | AMD Ryzen 5 1600"
                                            }}
                                        </li>
                                        <li>
                                            <strong>RAM:</strong>
                                            {{ selectedGame.recRAM || "16 GB" }}
                                        </li>
                                        <li>
                                            <strong>GPU:</strong>
                                            {{
                                                selectedGame.recGPU ||
                                                "NVIDIA GeForce GTX 1060 6GB | AMD Radeon RX 580 8GB"
                                            }}
                                        </li>
                                        <li>
                                            <strong>Storage:</strong>
                                            {{
                                                selectedGame.recStorage ||
                                                "60 GB SSD"
                                            }}
                                        </li>
                                    </ul>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import { LibraryManager } from "../../bindings/exhibition-launcher/library/index.js";

export default {
    name: "LibraryPage",

    data() {
        return {
            games: [],
            coverUrls: {},
            screenshotUrls: {},
            artworkUrls: {},
            currentPage: "library",
            selectedGame: null,
            currentSlide: 0,
        };
    },
    methods: {
        async loadCoverImages(game) {
            this.coverUrls[game.igdb_id] = await LibraryManager.GetCoverURL(
                game.cover_filename,
                game.cover_url,
            );
        },
        async loadArtworkImages(game) {
            let list = [];
            const artworks = await LibraryManager.GetAllImageURLs(
                game.artwork_filenames,
                game.artwork_url_list,
            );
            list.push(...artworks);
            this.artworkUrls[game.igdb_id] = list;
        },
        async loadScreenshotImages(game) {
            let list = [];
            const screens = await LibraryManager.GetAllImageURLs(
                game.screenshot_filenames,
                game.screenshot_url_list,
            );
            list.push(...screens);
            this.screenshotUrls[game.igdb_id] = list;
        },
        getGameImages() {
            let images = [];
            let game = this.selectedGame;
            // Use artwork images if available
            images.push(...this.artworkUrls[game.igdb_id]);
            images.push(...this.screenshotUrls[game.igdb_id]);

            // Use cover as fallback if no images available
            if (images.length === 0 && this.selectedGame?.cover_url) {
                images = this.coverUrls[game.igdb_id];
            }
            // At least 1 carousel image
            if (images.length === 0) {
                images = [
                    "https://via.placeholder.com/1920x1080/222222/555555?text=No+Images+Available",
                ];
            }

            return images;
        },

        nextSlide() {
            const totalSlides = this.getGameImages().length;
            this.currentSlide = (this.currentSlide + 1) % totalSlides;
        },

        prevSlide() {
            const totalSlides = this.getGameImages().length;
            this.currentSlide =
                (this.currentSlide - 1 + totalSlides) % totalSlides;
        },

        goToSlide(index) {
            this.currentSlide = index;
        },

        async addGame() {
            LibraryManager.AddToLibrary(119277, true)
                .catch((err) => {
                    console.warn(err);
                    return;
                })
                .then((game) => {
                    this.games.push(game);
                });
        },

        openGamePage(game) {
            this.selectedGame = game;
            this.currentPage = "game";
            this.currentSlide = 0; // reset carousel to first slide
        },

        launchGame(game) {
            LibraryManager.StartApp(game.igdb_id).catch((err) => {
                console.log(err);
            });
        },

        toggleFavorite() {
            if (this.selectedGame) {
                this.selectedGame.isFavorite = !this.selectedGame.isFavorite;
                console.log(
                    `Favorite status for game ${this.selectedGame.igdb_id}: ${this.selectedGame.isFavorite}`,
                );
            }
        },

        returnToLibrary() {
            this.currentPage = "library";
        },
    },
    async mounted() {
        const amountOfGames = await LibraryManager.GetAmountOfGames();

        console.log(amountOfGames);
        const portion = 100;

        for (let i = 0; i < amountOfGames; i += portion) {
            let games = await LibraryManager.GetRangeGame(portion, i);
            for (let j = 0; j < games.length; j++) {
                let game = games[j];
                await this.loadCoverImages(game);
                await this.loadArtworkImages(game);
                await this.loadScreenshotImages(game);
                this.games.push(game);
            }
        }

        // auto carousel each 5 seconds
        setInterval(() => {
            if (this.currentPage === "game") {
                this.nextSlide();
            }
        }, 5000);
    },
};
</script>

<style scoped>
.page {
    padding: 0;
    height: auto;
}

/* Library settings styles */
.Library-settings-container {
    width: 100%;
    height: 90px;
    display: flex;
    justify-content: space-between;
    align-items: center;
    position: sticky;
    top: 50px;
    padding: 50px;
    backdrop-filter: blur(15px);
    background-color: rgba(25, 25, 25, 0.6);
    z-index: 1;
    margin-bottom: 10px;
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
    border-radius: 20px;
    background-color: rgba(255, 255, 255, 0.1);
    border: 1px solid rgba(255, 255, 255, 0.1);
    transition: all 0.2s ease;
}

.Library-settings-wrapper button:hover {
    background-color: rgba(255, 255, 255, 0.2);
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
    padding: 5px 15px;
    position: relative;
    transition:
        color 0.3s ease,
        background-color 0.3s ease;
    border-radius: 20px;
}

.Library-favorites-container button:hover {
    background-color: rgba(255, 255, 255, 0.1);
}

.Library-favorites-container button.active {
    color: var(--accent-color);
    background-color: rgba(0, 0, 0, 0) !important;
}

.active-indicator-horizontal {
    position: absolute;
    bottom: -5px;
    left: 50%;
    transform: translateX(-50%);
    width: 0;
    height: 3px;
    background-color: var(--accent-color);
    border-radius: 2px;
    opacity: 0;
    transition:
        width 0.3s ease,
        opacity 0.3s ease;
}

.Library-favorites-container button.active .active-indicator-horizontal {
    width: 60%;
    opacity: 1;
}

.add-game-to-library-wrapper {
    display: flex;
    justify-content: center;
    align-items: center;
    gap: 10px;
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
    padding: 8px;
    border-radius: 50%;
    transition: background-color 0.2s ease;
}

#refresh-page:hover {
    background-color: rgba(255, 255, 255, 0.1);
}

.add-game-to-library-wrapper button.game-add-button {
    outline: none;
    border: none;
    padding: 10px 20px;
    color: var(--secondary-text-color);
    border-radius: 20px;
    background-color: rgba(255, 255, 255, 0.1);
    border: 1px solid rgba(255, 255, 255, 0.1);
    cursor: pointer;
    transition: all 0.2s ease;
    display: flex;
    gap: 10px;
    align-items: center;
}

.add-game-to-library-wrapper button:hover {
    background-color: rgba(255, 255, 255, 0.2);
    color: var(--text-color);
}

/* Game library container -*/
.game-library-container {
    width: 100%;
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(170px, 2fr));
    gap: 20px;
    padding-top: 20px;
    margin-bottom: 50px;
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
    transition:
        transform 0.3s ease,
        box-shadow 0.3s ease;
    background-position: center;
    background-repeat: no-repeat;
    background-size: cover;
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
}

.game-library-game-box:hover {
    transform: translateY(-5px);
    box-shadow: 0 8px 16px rgba(0, 0, 0, 0.3);
}

.game-box-info {
    width: 100%;
    background-color: rgba(0, 0, 0, 0.7);
    backdrop-filter: blur(10px);
    color: var(--text-color);
    font-size: 13px;
    position: absolute;
    bottom: -100px;
    padding: 12px 15px;
    transition: all 0.3s ease;
    display: flex;
    align-items: center;
    justify-content: space-between;
    border-bottom-left-radius: 15px;
    border-bottom-right-radius: 15px;
}

.game-box-info h1 {
    font-size: 14px;
    margin: 0;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    max-width: 80%;
}

.game-library-game-box:hover .game-box-info {
    bottom: 0;
}

.game-box-info .text-container {
    display: flex;
    flex-direction: column;
}

.game-box-info i {
    transition:
        color 0.3s ease,
        transform 0.3s ease;
    color: var(--secondary-text-color);
    font-size: 18px;
}

.game-box-info button {
    border: none;
    background: none;
    cursor: pointer;
    width: 30px;
    height: 30px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 50%;
    transition: background-color 0.2s ease;
}

.game-box-info button:hover {
    background-color: rgba(255, 255, 255, 0.1);
}

.game-box-info button:hover i {
    color: var(--text-color);
}

.game-page-new {
    width: 100%;
    min-height: 100vh;
    position: relative;
}

.back-button {
    position: fixed;
    top: 60px;
    margin-left: 10px;
    background: rgba(0, 0, 0, 0.6);
    backdrop-filter: blur(5px);
    border: none;
    color: white;
    cursor: pointer;
    padding: 10px;
    border-radius: 20px;
    transition: all 0.2s ease;
    z-index: 10;
    font-size: 10px;
    display: flex;
    align-items: center;
}

.back-button:hover {
    background: rgba(0, 0, 0, 0.8);
}

/* Carousel Styles */
.carousel-container {
    position: relative;
    width: 100%;
    height: 600px;
    overflow: hidden;
}

.carousel-slides {
    display: flex;
    height: 100%;
    transition: transform 0.5s ease-in-out;
}

.carousel-slide {
    min-width: 100%;
    height: 100%;
    background-size: cover;
    background-position: center;
    position: relative;
}

.banner-overlay {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background: linear-gradient(
        to bottom,
        rgba(0, 0, 0, 0.1) 0%,
        rgba(0, 0, 0, 0.3) 50%,
        rgba(0, 0, 0, 0.8) 90%
    );
}

.carousel-controls {
    position: absolute;
    bottom: 30px;
    left: 0;
    width: 100%;
    display: flex;
    justify-content: right;
    align-items: center;
    gap: 30px;
    z-index: 5;
    padding: 0 5%;
}

.carousel-btn {
    background: rgba(0, 0, 0, 0.5);
    border: none;
    color: white;
    width: 40px;
    height: 40px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    transition: background-color 0.2s ease;
}

.carousel-btn:hover {
    background: rgba(0, 0, 0, 0.8);
}

.carousel-indicators {
    display: flex;
    gap: 8px;
}

.carousel-indicators span {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    background-color: rgba(255, 255, 255, 0.5);
    cursor: pointer;
    transition: all 0.2s ease;
}

.carousel-indicators span.active {
    background-color: white;
    transform: scale(1.2);
}

/* Game Info Overlay */
.game-info-overlay {
    position: absolute;
    bottom: 0;
    left: 0;
    width: 100%;
    padding: 0 5%;
    z-index: 2;
}

.game-branding {
    margin-bottom: 15px;
}

.game-title {
    font-size: 3rem;
    font-weight: bold;
    color: white;
    margin: 0;
    text-shadow: 1px 1px 3px rgba(0, 0, 0, 0.5);
}

.game-description {
    max-width: 600px;
    margin-bottom: 20px;
}

.game-description p {
    color: rgba(255, 255, 255, 0.9);
    font-size: 1rem;
    line-height: 1.5;
}

.game-stats-container {
    display: flex;
    margin-bottom: 30px;
}

.stat-box {
    background-color: rgba(0, 0, 0, 0.6);
    backdrop-filter: blur(5px);
    padding: 15px 30px;
    text-align: center;
    border-radius: 10px;
    margin-right: 15px;
}

.stat-value {
    color: white;
    font-size: 1.25rem;
    font-weight: bold;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
}

.stat-label {
    color: rgba(255, 255, 255, 0.7);
    font-size: 0.75rem;
    text-transform: uppercase;
    margin-top: 4px;
}

.game-actions {
    display: flex;
    gap: 15px;
    margin-bottom: 30px;
}

.add-library-btn,
.download-btn {
    padding: 12px 25px;
    border-radius: 25px;
    border: none;
    font-size: 0.9rem;
    font-weight: 500;
    cursor: pointer;
    display: flex;
    align-items: center;
    gap: 10px;
    transition: all 0.3s ease;
}

.add-library-btn {
    background-color: rgba(255, 255, 255, 0.2);
    backdrop-filter: blur(5px);
    color: white;
}

.download-btn {
    background-color: var(--accent-color, #3a86ff);
    color: white;
}

.add-library-btn:hover {
    background-color: rgba(255, 255, 255, 0.3);
}

.download-btn:hover {
    filter: brightness(1.1);
}

/* Game Details Section */
.game-details-section {
    padding: 30px 5%;
}

.game-details-container {
    display: flex;
    gap: 50px;
    flex-wrap: wrap;
}

.details-column {
    flex: 1;
    min-width: 300px;
}

.details-section {
    margin-bottom: 30px;
}

.details-section h3 {
    font-size: 1.25rem;
    margin-bottom: 15px;
    color: var(--text-color, white);
    font-weight: 600;
}

.details-section p {
    color: var(--secondary-text-color, rgba(255, 255, 255, 0.7));
    line-height: 1.6;
}

.tags-container {
    display: flex;
    flex-wrap: wrap;
    gap: 10px;
}

.tag {
    background-color: rgba(255, 255, 255, 0.1);
    color: var(--secondary-text-color, rgba(255, 255, 255, 0.8));
    padding: 8px 15px;
    border-radius: 20px;
    font-size: 0.875rem;
}

.system-reqs {
    display: flex;
    flex-wrap: wrap;
    gap: 30px;
    background-color: rgba(255, 255, 255, 0.05);
    border-radius: 15px;
    padding: 20px;
}

.req-section {
    flex: 1;
    min-width: 250px;
}

.req-section h4 {
    color: var(--text-color, white);
    margin-bottom: 15px;
    font-size: 1rem;
    font-weight: 500;
}

.req-section ul {
    list-style: none;
    padding: 0;
    margin: 0;
}

.req-section li {
    color: var(--secondary-text-color, rgba(255, 255, 255, 0.7));
    margin-bottom: 8px;
    font-size: 0.9rem;
}

.req-section strong {
    color: var(--text-color, white);
}

/* Media queries for better responsive design */
@media (max-width: 768px) {
    .carousel-container {
        height: 450px;
    }

    .game-title {
        font-size: 2rem;
    }

    .carousel-controls {
        bottom: 70px;
    }

    .system-reqs {
        flex-direction: column;
        gap: 20px;
    }
}

@media (max-width: 480px) {
    .carousel-container {
        height: 350px;
    }

    .game-title {
        font-size: 1.5rem;
    }

    .game-actions {
        flex-direction: column;
        width: 100%;
    }

    .add-library-btn,
    .download-btn {
        width: 100%;
        justify-content: center;
    }

    .stat-box {
        padding: 10px 15px;
    }
}
</style>
