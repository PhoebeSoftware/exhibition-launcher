<template>
    <div class="page">
        <div>
            <h1>Explore</h1>
            <p>Discover new content here.</p>
            <button @click="addToQueue()">Add torrent by magnet</button>
            <button @click="startDownloads()">Start downloads / q</button>
            <button @click="pauseResumeDownloads()">Pause / Resume download</button>
            <p>Download progress {{ progress.toFixed(2) }}%</p>
            <p>Time passed: {{ timePassed }}</p>
            <p>Paused: {{ paused }}</p>
        </div>
        <div>
            <form @submit.prevent="addGames">
                <input type="text" v-model="name" placeholder="Search to add game">
                <button type="submit">Add game(s)</button>
            </form>
        </div>
    </div>
</template>


<script setup>
import {
    AddRealDebridDownloadToQueue,
    StartDownloads
} from "../../bindings/exhibition-launcher/exhibition_queue/queue.js";
import {Settings} from "../../bindings/exhibition-launcher/utils/json_utils/json_models/index.js";

async function addToQueue(magnetLink) {

    let settings = await Settings.GetSettings();

    if (!settings.real_debrid_settings.use_real_debrid) {
        return;
    }

    await AddRealDebridDownloadToQueue(magnetLink);
}


function startDownloads() {
    StartDownloads().catch((err) => {
        console.log(err);
    });
}
</script>

<script>
import {Events} from "@wailsio/runtime";
import {Queue} from "../../bindings/exhibition-launcher/exhibition_queue/index.js";
import {LibraryManager} from "../../bindings/exhibition-launcher/library/index.js";
import {ProxyClient} from "../../bindings/exhibition-launcher/proxy_client/index.js";

export default {
    name: "ExplorePage",
    async mounted() {
        Events.On("download_progress", this.updateDownloadProgress);
        this.paused = await Queue.GetPaused();
    },
    methods: {
        async addGames() {
            let games = await ProxyClient.GetMetadataByName(this.name)
            for (let game of games) {
               LibraryManager.AddToLibrary(game.id)
            }
        },
        updateDownloadProgress(event) {
            const progressData = event.data[0] || event.data["@"];

            if (progressData) {
                this.progress = progressData.percent || 0;
                this.timePassed = progressData.timePassed || "";
                console.log("Current progress:", this.progress);

                console.log("Downloaded bytes:", progressData.downloadedBytes);
                console.log("Total bytes:", progressData.totalBytes);
                console.log("Time passed:", progressData.timePassed);
            } else {
                console.error("Progress data not found in expected format");
            }
        },
        async pauseResumeDownloads() {
            let pauseValue = !(await Queue.GetPaused());
            Queue.SetPaused(pauseValue);
            this.paused = pauseValue;
            console.log(pauseValue);
        },

    },
    data() {
        return {
            progress: 0,
            paused: false,
            timePassed: "",
            name: '',
        };
    },
};
</script>

<style scoped>
.page {
    padding: 20px;
}

h1 {
    margin-bottom: 15px;
}
</style>
