# Company of Heroes 3 Replay Manager

[![Build Status](https://github.com/SavageCore/coh3-replay-manager-go/actions/workflows/ci.yml/badge.svg?style=for-the-badge)](https://github.com/SavageCore/coh3-replay-manager-go/actions/workflows/ci.yml)
[![GitHub all releases](https://img.shields.io/github/downloads/SavageCore/coh3-replay-manager-go/total?style=flat)](https://github.com/SavageCore/coh3-replay-manager-go/releases/latest)

> Tray app that automatically saves your replays and allows direct download and
> play from [cohdb.com](https://cohdb.com/)

So, how's this all work then?

The app watches the directory
`%USERPROFILE%\Documents\My Games\Company of Heroes 3\playback` which is where
the replay of your last game played is stored as `temp.rec` and
`temp_campaign.rec` it then automatically renames these files to saved-replay
followed by Unix time.

Playing a replay from [cohdb.com](https://cohdb.com/) is also supported. The app
will download the replay and then launch the game with a single click!

Needs a decent logo and a better name probably. I'm open to suggestions and PRs.

# Running 🏃

Download the
[latest release](https://github.com/SavageCore/coh3-replay-manager-go/releases/latest)
and run it 🚀

⚠️ For now, you'll also need to install a User script to get the play button to
display on [cohdb.com](https://cohdb.com/). ⚠️

These are the steps:

1. Install [Greasemonkey](https://www.greasespot.net/) or
   [Violentmonkey](https://violentmonkey.github.io/get-it/) for your browser.
2. Click
   [here](https://cdn.jsdelivr.net/gh/SavageCore/coh3-replay-manager-go/userscript/src/coh3-replay-manager-go.user.js)
   to install the User script.
3. You should now see a play button next to the download button on
   [cohdb.com](https://cohdb.com/) 🎉

# Development 🖥️

Install the dependencies:

1. [Go](https://go.dev/doc/install)
1. [Task](https://taskfile.dev/)
1. [Air](https://github.com/cosmtrek/air)
1. [go-bindata](https://github.com/go-bindata/go-bindata)
1. [rsrc](https://github.com/akavel/rsrc)
1. [staticcheck](https://staticcheck.io/docs/getting-started/)

Clone the repo:

1. `git clone https://github.com/SavageCore/coh3-replay-manager-go`
1. `cd coh3-replay-manager-go`

Run the app with live reloading:

1. Run `task dev`

# Building 🚧

1. Run `task build` and the app will be built to `build/bin/coh3-replay-manager-go.exe`

# Releasing 🚀

1. Increment `CurrentVersion` in `main.go`
2. Push a git tag with the new version number
3. GitHub Action will build and release the app 🎉
