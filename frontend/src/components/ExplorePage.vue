<template>
    <div class="page">
        <h1>Explore</h1>
        <p>Discover new content here.</p>
        <button @click="addToQueue()">Add torrent by magnet</button>
        <button @click="startDownloads()">Start downloads / q</button>
        <button @click="pauseResumeDownloads()">Pause / Resume download</button>
        <p>Download progress {{ progress.toFixed(2) }}%</p>
        <p>Time passed: {{ timePassed }}</p>
        <p>Paused: {{ paused }}</p>
    </div>
</template>

<script setup>
import { Settings } from "../../bindings/exhibition-launcher/utils/jsonUtils/jsonModels/index.js";
import { RealDebridClient } from "../../bindings/exhibition-launcher/torrent/realdebrid/index.js";

async function addToQueue(magnetLink) {
    let magnetLinkHollowKnight =
        "magnet:?xt=urn:btih:D738F320446AEB504C80904F670B0615D04D5C6C&dn=Hollow+Knight+%28v1.5.68.11808+%2B+2+Bonus+OSTs%2C+MULTi10%29+%5BFitGirl+Repack%2C+Selective+Download+-+from+814+MB%5D&tr=udp%3A%2F%2F46.148.18.250%3A2710&tr=udp%3A%2F%2Fopentor.org%3A2710&tr=udp%3A%2F%2Ftracker.coppersurfer.tk%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337%2Fannounce&tr=http%3A%2F%2Ftracker.dler.org%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.leechers-paradise.org%3A6969%2Fannounce&tr=udp%3A%2F%2F9.rarbg.me%3A2730%2Fannounce&tr=udp%3A%2F%2F9.rarbg.to%3A2770%2Fannounce&tr=udp%3A%2F%2Ftracker.pirateparty.gr%3A6969%2Fannounce&tr=http%3A%2F%2Fretracker.local%2Fannounce&tr=http%3A%2F%2Fretracker.ip.ncnet.ru%2Fannounce&tr=udp%3A%2F%2Fexodus.desync.com%3A6969%2Fannounce&tr=udp%3A%2F%2Fipv4.tracker.harry.lu%3A80%2Fannounce&tr=udp%3A%2F%2Fopen.stealth.si%3A80%2Fannounce&tr=udp%3A%2F%2Ftracker.zer0day.to%3A1337%2Fannounce&tr=udp%3A%2F%2Ftracker.leechers-paradise.org%3A6969%2Fannounce&tr=udp%3A%2F%2Fcoppersurfer.tk%3A6969%2Fannounce";
    let idkrandomapp =
        "magnet:?xt=urn:btih:EEEF75F8C7AD79818C54C618E1A7937CD76B59C4&dn=Sony+Vegas+Pro+v11.0.510+64+bit+%28patch+keygen+DI%29+%5BChingLiu%5D&tr=http%3A%2F%2Fpow7.com%2Fannounce&tr=http%3A%2F%2Fpubt.net%3A2710%2Fannounce&tr=http%3A%2F%2Ft1.pow7.com%2Fannounce&tr=http%3A%2F%2Ftracker.torrentbay.to%3A6969%2Fannounce&tr=http%3A%2F%2Ftracker.torrent.to%3A2710%2Fannounce&tr=http%3A%2F%2Ftracker.publicbt.com%2Fannounce&tr=udp%3A%2F%2Ftracker.1337x.org%3A80%2Fannounce&tr=udp%3A%2F%2Ftracker.openbittorrent.com%3A80%2Fannounce&tr=udp%3A%2F%2Ftracker.istole.it%3A80%2Fannounce&tr=http%3A%2F%2Ftracker.publicbt.com%3A80%2Fannounce&tr=http%3A%2F%2Fa.tracker.prq.to%2Fannounce&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337%2Fannounce&tr=http%3A%2F%2Ftracker.openbittorrent.com%3A80%2Fannounce&tr=udp%3A%2F%2Fopentracker.i2p.rocks%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.internetwarriors.net%3A1337%2Fannounce&tr=udp%3A%2F%2Ftracker.leechers-paradise.org%3A6969%2Fannounce&tr=udp%3A%2F%2Fcoppersurfer.tk%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.zer0day.to%3A1337%2Fannounce";

    let settings = await Settings.GetSettings();

    if (!settings.real_debrid_settings.use_real_debrid) {
        return;
    }

    await Queue.AddRealDebridDownloadToQueue(magnetLinkHollowKnight);
}

function startDownloads() {
    Queue.StartDownloads().catch((err) => {
        console.log(err);
    });
}
</script>

<script>
import { Events } from "@wailsio/runtime";
import { RealDebridClient } from "../../bindings/exhibition-launcher/torrent/realdebrid/index.js";
import { Queue } from "../../bindings/exhibition-launcher/exhibitionQueue/index.js";

export default {
    name: "ExplorePage",
    async mounted() {
        Events.On("download_progress", this.updateDownloadProgress);
        this.paused = await Queue.GetPaused();
    },
    methods: {
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
