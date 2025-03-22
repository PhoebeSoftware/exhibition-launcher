<template>
    <div class="page">
        <h1>Explore</h1>
        <p>Discover new content here.</p>
        <button @click="addTorrentByMagnet()">Add torrent by magnet</button>
    </div>

</template>

<script setup>

import {RealDebridClient} from "../../bindings/derpy-launcher072/torrent/realdebrid/index.js";
import {SettingsManager} from "../../bindings/derpy-launcher072/utils/settings/index.js";

async function addTorrentByMagnet() {
    // RealDebridClient.GetTrafic().then((result) => {
    //     console.log(result)
    // });
    // RealDebridClient.GetUser().then((result) => {
    //     console.log(result)
    // });
    let id;
    let magnetLink = "magnet:?xt=urn:btih:BED99329D89D10663336AE1D76ACAFF5EFE27EAA&dn=Bleach%3A+Rebirth+of+Souls+-+Ultimate+Edition+%28v1.03+%2B+All+DLCs+%2B+MULTi15%29+%28Fast+Install%29+%5BDODI+Repack%5D&tr=udp%3A%2F%2F9.rarbg.to%3A2870%2Fannounce&tr=udp%3A%2F%2Feddie4.nl%3A6969%2Fannounce&tr=udp%3A%2F%2Fthetracker.org%3A80%2Fannounce&tr=udp%3A%2F%2Ftracker.coppersurfer.tk%3A6969&tr=udp%3A%2F%2Ftracker.leechers-paradise.org%3A6969&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337&tr=udp%3A%2F%2Ftracker.zer0day.to%3A1337&tr=udp%3A%2F%2F9.rarbg.com%3A2710%2Fannounce&tr=udp%3A%2F%2Fexplodie.org%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker1.wasabii.com.tw%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.torrent.eu.org%3A451%2Fannounce&tr=udp%3A%2F%2Ftracker.cypherpunks.ru%3A6969%2Fannounce&tr=udp%3A%2F%2Fp4p.arenabg.com%3A1337%2Fannounce&tr=udp%3A%2F%2Ftracker.opentrackr.org%3A1337%2Fannounce&tr=http%3A%2F%2Ftracker.openbittorrent.com%3A80%2Fannounce&tr=udp%3A%2F%2Fopentracker.i2p.rocks%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.internetwarriors.net%3A1337%2Fannounce&tr=udp%3A%2F%2Ftracker.leechers-paradise.org%3A6969%2Fannounce&tr=udp%3A%2F%2Fcoppersurfer.tk%3A6969%2Fannounce&tr=udp%3A%2F%2Ftracker.zer0day.to%3A1337%2Fannounce"

    let settings = await SettingsManager.GetSettings()

    await RealDebridClient.DownloadByMagnet(magnetLink, settings);

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