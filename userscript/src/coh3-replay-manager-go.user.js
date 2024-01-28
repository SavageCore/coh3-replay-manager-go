// ==UserScript==
// @name         coh3-replay-manager-go-helper
// @namespace    https://savagecore.uk
// @version      0.1.0
// @description  Add play button to cohdb.com to invoke coh3-replay-manager-go to download and play the replay
// @author       SavageCore
// @include      https://cohdb.com/
// @include      https://cohdb.com/replays
// @include      https://cohdb.com/replays?page=*
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
    // Example download button:
    // #replays > div > div:nth-child(1) > div:nth-child(8) > form
    // form has action="/replays/1391/file" and class button_to
    // The row that contains the button has class of list-group-item
    const downloadButtons = document.querySelectorAll('form.button_to[action^="/replays/"] > button.btn.btn-outline-primary');
    for (const downloadButton of downloadButtons) {
        const row = downloadButton.parentElement.parentElement.parentElement;
        const versionElement = row.querySelector('small > span:last-child')
        const gameVersion = versionElement ? versionElement.textContent.trim() : '0';
        const replayId = downloadButton.parentElement.action.split('/')[4];

        const playButton = downloadButton.cloneNode(true);

        playButton.querySelector('i').classList.remove('fa-download');
        playButton.querySelector('i').classList.add('fa-play');
        playButton.classList.remove('btn-outline-primary');
        playButton.classList.add('btn-outline-success');
        playButton.querySelector('i').nextSibling.textContent = ' ';
        playButton.style.height = `${downloadButton.clientHeight + 2}px`;
        playButton.title = 'Play with coh3-replay-manager-go';

        playButton.addEventListener('click', e => {
            e.preventDefault();
            console.log('Play clicked');
            window.location.href = `coh3-replay-manager-go://play/${replayId}/v/${gameVersion}`;
        });

        downloadButton.addEventListener('click', e => {
            e.preventDefault();
            console.log('Download clicked');
            window.location.href = `coh3-replay-manager-go://download/${replayId}/v/${gameVersion}`;
        });

        downloadButton.parentNode.parentNode.insertBefore(playButton, downloadButton.nextSibling);
    }
}

main();
