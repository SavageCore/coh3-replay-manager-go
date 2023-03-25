// ==UserScript==
// @name         coh3-replay-manager-go-helper
// @namespace    https://savagecore.uk
// @version      0.1.0
// @description  Add play button to cohdb.com to invoke coh3-replay-manager-go to download and play the replay
// @author       SavageCore
// @include      https://cohdb.com/
// @include      https://cohdb.com/?page=*
// @require      https://greasemonkey.github.io/gm4-polyfill/gm4-polyfill.js
// @run-at       document-idle
// ==/UserScript==

const main = async () => {
    addPlayButtons();

    const html = document.documentElement;

    // Create a new MutationObserver
    const observer = new MutationObserver(function (mutations) {
        mutations.forEach(function (mutation) {
            // Check if the aria-busy and data-turbo-preview attributes have changed which indicates the page has been updated
            // We then add the play buttons if they are not already present
            if (mutation.attributeName === 'aria-busy' || mutation.attributeName === 'data-turbo-preview') {
                if (!html.hasAttribute('aria-busy') && !html.hasAttribute('data-turbo-preview')) {
                    const playButtons = document.querySelectorAll('a[href^="coh3-replay-manager-go://play"]');
                    if (playButtons.length === 0) {
                        console.log('coh3-replay-manager-go-helper: Adding play buttons');
                        addPlayButtons();
                    }
                }
            }
        });
    });

    // Start observing changes to the html element
    observer.observe(html, { attributes: true });
}

const addPlayButtons = () => {
    const downloadButtons = document.querySelectorAll('a[href$="/download"]');
    for (const button of downloadButtons) {
        const versionElement = button.parentElement.querySelector('small > span');
        const gameVersion = versionElement ? versionElement.textContent.replace('v', '') : '0';
        const replayId = button.href.split('/')[4];

        const playButton = button.cloneNode(true);

        playButton.href = `coh3-replay-manager-go://play/${replayId}/v/${gameVersion}`;
        playButton.querySelector('i').classList.remove('fa-download');
        playButton.querySelector('i').classList.add('fa-play');
        playButton.classList.remove('btn-outline-primary');
        playButton.classList.add('btn-outline-success');
        playButton.querySelector('i').nextSibling.textContent = ' ';
        playButton.style.height = `${button.clientHeight + 2}px`;
        playButton.title = 'Play with coh3-replay-manager-go';

        button.parentNode.insertBefore(playButton, button.nextSibling);
    }
}

main();
