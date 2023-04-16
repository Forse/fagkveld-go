import "./style.css";
import "reveal.js/plugin/highlight/monokai.css"
import "reveal.js/dist/reveal.css";
// see available themes in the
// node_modules/reveal.js/dist/theme
//  beige, black, blood, league, moon, night, serif, simple, ...
import "reveal.js/dist/theme/black.css";
import Reveal from "reveal.js";
import Markdown from "reveal.js/plugin/markdown/markdown";
import Highlight from "reveal.js/plugin/highlight/highlight";

const deck = new Reveal();
deck.initialize({
    plugins: [Markdown, Highlight],
    hash: true,
    slideNumber: true,
    // width: 3840,
    // height: 2160,
    // center: false,
    progress: true,
});
