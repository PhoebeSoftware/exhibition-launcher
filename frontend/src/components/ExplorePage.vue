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
    let magnetLink = "magnet:?xt=urn:btih:ETIZG6JR5OHURVXIVZIPYMPQ3EZO3AGA&dn=kali-linux-2025.1a-installer-arm64.iso&xl=3758403584&tr=http%3A%2F%2Ftracker.kali.org%3A6969%2Fannounce"

    let settings = await Settings.GetSettings()
    try {
        await RealDebridClient.DownloadByMagnet(magnetLink, settings.download_path);

    } catch (e) {
        console.log(e)
    }

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