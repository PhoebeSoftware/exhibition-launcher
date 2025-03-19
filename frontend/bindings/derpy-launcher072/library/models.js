// @ts-check
// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore: Unused imports
import {Create as $Create} from "@wailsio/runtime";

export class Game {
    /**
     * Creates a new Game instance.
     * @param {Partial<Game>} [$$source = {}] - The source object to create the Game.
     */
    constructor($$source = {}) {
        if (!("igdb_id" in $$source)) {
            /**
             * @member
             * @type {number}
             */
            this["igdb_id"] = 0;
        }
        if (!("playtime" in $$source)) {
            /**
             * @member
             * @type {number}
             */
            this["playtime"] = 0;
        }
        if (!("achievments" in $$source)) {
            /**
             * @member
             * @type {number[]}
             */
            this["achievments"] = [];
        }
        if (!("executable" in $$source)) {
            /**
             * @member
             * @type {string}
             */
            this["executable"] = "";
        }
        if (!("running" in $$source)) {
            /**
             * @member
             * @type {boolean}
             */
            this["running"] = false;
        }
        if (!("favorite" in $$source)) {
            /**
             * @member
             * @type {boolean}
             */
            this["favorite"] = false;
        }

        Object.assign(this, $$source);
    }

    /**
     * Creates a new Game instance from a string or object.
     * @param {any} [$$source = {}]
     * @returns {Game}
     */
    static createFrom($$source = {}) {
        const $$createField2_0 = $$createType0;
        let $$parsedSource = typeof $$source === 'string' ? JSON.parse($$source) : $$source;
        if ("achievments" in $$parsedSource) {
            $$parsedSource["achievments"] = $$createField2_0($$parsedSource["achievments"]);
        }
        return new Game(/** @type {Partial<Game>} */($$parsedSource));
    }
}

// Private type creation functions
const $$createType0 = $Create.Array($Create.Any);
