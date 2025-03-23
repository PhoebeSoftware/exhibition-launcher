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
    let magnetLink = "magnet:?xt=urn:btih:7CA678ADFD25AA3697E8E7FB5C1ABA4CDB70F843&dn=Shingeki+no+Kyojin+%28Attack+on+Titan%29+%28Season+4+Part+2%29+%5B1080p%5D%5BHEVC+x265+10bit%5D%5BMulti-Subs%5D+-+Judas&tr=http%3A%2F%2Fnyaa.tracker.wf%3A7777%2Fannounce&tr=udp%3A%2F%2Fopen.stealth.si%3A80%2Fannounce&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337%2Fannounce&tr=udp%3A%2F%2Fexodus.desync.com%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.torrent.eu.org%3A451%2Fannounce&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337%2Fannounce&tr=http%3A%2F%2Ftracker.openbittorrent.com%3A80%2Fannounce&tr=udp%3A%2F%2Fopentracker.i2p.rocks%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.internetwarriors.net%3A1337%2Fannounce&tr=udp%3A%2F%2Ftracker.leechers-paradise.org%3A6969%2Fannounce&tr=udp%3A%2F%2Fcoppersurfer.tk%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.zer0day.to%3A1337%2Fannounce"

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