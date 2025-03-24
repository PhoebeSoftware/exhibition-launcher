<template>
    <div class="page">
        <h1>Explore</h1>
        <p>Discover new content here.</p>
        <button @click="addTorrentByMagnet()">Add torrent by magnet</button>
    </div>

</template>

<script setup>

import {RealDebridClient} from "../../bindings/derpy-launcher072/torrent/realdebrid/index.js";
import {GetSettings} from "../../bindings/derpy-launcher072/utils/jsonUtils/jsonModels/settings.js";
import {JsonManager} from "../../bindings/derpy-launcher072/utils/jsonUtils/index.js";
import {Settings} from "../../bindings/derpy-launcher072/utils/jsonUtils/jsonModels/index.js";

async function addTorrentByMagnet() {
    // RealDebridClient.GetTrafic().then((result) => {
    //     console.log(result)
    // });
    // RealDebridClient.GetUser().then((result) => {
    //     console.log(result)
    // });
    let id;
    let magnetLink = "magnet:?xt=urn:btih:D738F320446AEB504C80904F670B0615D04D5C6C&dn=Hollow+Knight+%28v1.5.68.11808+%2B+2+Bonus+OSTs%2C+MULTi10%29+%5BFitGirl+Repack%2C+Selective+Download+-+from+814+MB%5D&tr=udp%3A%2F%2F46.148.18.250%3A2710&tr=udp%3A%2F%2Fopentor.org%3A2710&tr=udp%3A%2F%2Ftracker.coppersurfer.tk%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337%2Fannounce&tr=http%3A%2F%2Ftracker.dler.org%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.leechers-paradise.org%3A6969%2Fannounce&tr=udp%3A%2F%2F9.rarbg.me%3A2730%2Fannounce&tr=udp%3A%2F%2F9.rarbg.to%3A2770%2Fannounce&tr=udp%3A%2F%2Ftracker.pirateparty.gr%3A6969%2Fannounce&tr=http%3A%2F%2Fretracker.local%2Fannounce&tr=http%3A%2F%2Fretracker.ip.ncnet.ru%2Fannounce&tr=udp%3A%2F%2Fexodus.desync.com%3A6969%2Fannounce&tr=udp%3A%2F%2Fipv4.tracker.harry.lu%3A80%2Fannounce&tr=udp%3A%2F%2Fopen.stealth.si%3A80%2Fannounce&tr=udp%3A%2F%2Ftracker.zer0day.to%3A1337%2Fannounce&tr=udp%3A%2F%2Ftracker.leechers-paradise.org%3A6969%2Fannounce&tr=udp%3A%2F%2Fcoppersurfer.tk%3A6969%2Fannounce"

    let settings = await Settings.GetSettings()
    RealDebridClient.DownloadByMagnet(magnetLink, settings.download_path).catch((err) => {
        console.log(err)
    });

    //await RealDebridClient.DownloadByMagnet(magnetLink, Settings.)


    /*    await RealDebridClient.AddTorrentByMagnet(magnetLink).then((result) => {
                id = result.id
                console.log(result)
            })

        console.log(id)

        let torrent;

        await RealDebridClient.GetTorrentInfoById(id).then((result) => {
            torrent = result
            console.log(result)
        });

        await RealDebridClient.SelectFiles(torrent);


        let torrentWithSelectedFiles;
        await RealDebridClient.GetTorrentInfoById(id).then((result) => {
            torrentWithSelectedFiles = result
            console.log(result)
        });

        let downloadLinks = []

        await torrentWithSelectedFiles.links.forEach((torrentWithSelectedFiles) => {
            RealDebridClient.UnrestrictLink(torrentWithSelectedFiles).then((unrestrictObject) => {
                downloadLinks.push(unrestrictObject)
            })
        })

        await downloadLinks.forEach((unrestrictObject) => {
        })*/

    // RealDebridClient.GetDownloads().then((result) => {
    //     console.log(result)
    // });
    // RealDebridClient.GetTorents().then((result) => {
    //     console.log(result)
    // });
}
</script>

<script>
export default {
    name: 'ExplorePage'
}

</script>

<style scoped>
.page {
    padding: 20px;
}

h1 {
    margin-bottom: 15px;
}
</style>