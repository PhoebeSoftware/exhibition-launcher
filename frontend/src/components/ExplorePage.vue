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
    let magnetLink = 'magnet:?xt=urn:btih:EEEF75F8C7AD79818C54C618E1A7937CD76B59C4&dn=Sony+Vegas+Pro+v11.0.510+64+bit+%28patch+keygen+DI%29+%5BChingLiu%5D&tr=http%3A%2F%2Fpow7.com%2Fannounce&tr=http%3A%2F%2Fpubt.net%3A2710%2Fannounce&tr=http%3A%2F%2Ft1.pow7.com%2Fannounce&tr=http%3A%2F%2Ftracker.torrentbay.to%3A6969%2Fannounce&tr=http%3A%2F%2Ftracker.torrent.to%3A2710%2Fannounce&tr=http%3A%2F%2Ftracker.publicbt.com%2Fannounce&tr=udp%3A%2F%2Ftracker.1337x.org%3A80%2Fannounce&tr=udp%3A%2F%2Ftracker.openbittorrent.com%3A80%2Fannounce&tr=udp%3A%2F%2Ftracker.istole.it%3A80%2Fannounce&tr=http%3A%2F%2Ftracker.publicbt.com%3A80%2Fannounce&tr=http%3A%2F%2Fa.tracker.prq.to%2Fannounce&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337%2Fannounce&tr=http%3A%2F%2Ftracker.openbittorrent.com%3A80%2Fannounce&tr=udp%3A%2F%2Fopentracker.i2p.rocks%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.internetwarriors.net%3A1337%2Fannounce&tr=udp%3A%2F%2Ftracker.leechers-paradise.org%3A6969%2Fannounce&tr=udp%3A%2F%2Fcoppersurfer.tk%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.zer0day.to%3A1337%2Fannounce'

    let settings = await Settings.GetSettings()
    try {
        RealDebridClient.DownloadByMagnet(magnetLink, settings.download_path).catch((err) => {
            console.log(err)
        });
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