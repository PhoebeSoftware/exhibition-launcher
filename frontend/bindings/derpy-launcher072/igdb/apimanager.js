// @ts-check
// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore: Unused imports
import {Call as $Call, Create as $Create} from "@wailsio/runtime";

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore: Unused imports
import * as $models from "./models.js";

/**
 * @param {number} coverID
 * @returns {Promise<string> & { cancel(): void }}
 */
export function GetCover(coverID) {
    let $resultPromise = /** @type {any} */($Call.ByID(2756693912, coverID));
    return $resultPromise;
}

/**
 * @param {number} id
 * @returns {Promise<$models.ApiGame> & { cancel(): void }}
 */
export function GetGameData(id) {
    let $resultPromise = /** @type {any} */($Call.ByID(511432535, id));
    let $typingPromise = /** @type {any} */($resultPromise.then(($result) => {
        return $$createType0($result);
    }));
    $typingPromise.cancel = $resultPromise.cancel.bind($resultPromise);
    return $typingPromise;
}

/**
 * @param {string} query
 * @returns {Promise<$models.ApiGame[]> & { cancel(): void }}
 */
export function GetGames(query) {
    let $resultPromise = /** @type {any} */($Call.ByID(195213204, query));
    let $typingPromise = /** @type {any} */($resultPromise.then(($result) => {
        return $$createType1($result);
    }));
    $typingPromise.cancel = $resultPromise.cancel.bind($resultPromise);
    return $typingPromise;
}

// Private type creation functions
const $$createType0 = $models.ApiGame.createFrom;
const $$createType1 = $Create.Array($$createType0);
