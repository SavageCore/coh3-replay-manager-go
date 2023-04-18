<script lang="ts">
  import {
    List,
    Play,
    PlayLocal,
    Remove,
    RemoveLocal,
  } from "../wailsjs/go/main/App.js";
  import "chota";
  import { Modal, Button, Row, Col, Icon, Card } from "svelte-chota";
  import { mdiDelete, mdiPlay } from "@mdi/js";

  // Import map icons
  import aerePerenniusIcon from "./assets/icons/maps/aere_perennius.webp";
  import campbellsConvoyIcon from "./assets/icons/maps/campbells_convoy.webp";
  import gardensIcon from "./assets/icons/maps/gardens.png";
  import crossroadsIcon from "./assets/icons/maps/gardens.png";
  import gazalaLandingGroundIcon from "./assets/icons/maps/gazala_landing_ground.webp";
  import laquilaIcon from "./assets/icons/maps/laquila.webp";
  import mignanoGap6pIcon from "./assets/icons/maps/mignano_gap_6p.png";
  import mignanoGap8pIcon from "./assets/icons/maps/mignano_gap.webp";
  import pachinoFarmlandsIcon from "./assets/icons/maps/pachino_farmlands.webp";
  import pachinoFarmlandsMkiiIcon from "./assets/icons/maps/pachino_farmlands_mkii.png";
  import roadToTunisIcon from "./assets/icons/maps/road_to_tunis.webp";
  import tarantoCoastlineIcon from "./assets/icons/maps/taranto_coastline.webp";
  import torrenteIcon from "./assets/icons/maps/torrente.webp";
  import twinBeachesIcon from "./assets/icons/maps/twin_beaches.webp";
  import winterLineIcon from "./assets/icons/maps/winter_line.webp";

  // Import faction icons
  import dakIcon from "./assets/icons/factions/dak.webp";
  import britishIcon from "./assets/icons/factions/british.webp";
  import wehrmachtIcon from "./assets/icons/factions/german.webp";
  import americanIcon from "./assets/icons/factions/american.webp";

  let replays;
  let localReplays;
  let localRemoveTarget;
  let removeTarget;

  let openRemove = false;
  let openLocalRemove = false;
  let open = false;
  let size;
  const modal_show = () => (openRemove = true);
  const modal_hide = () => (openRemove = false);
  const modal_show_local = () => (openLocalRemove = true);
  const modal_hide_local = () => (openLocalRemove = false);

  if (
    window.matchMedia &&
    window.matchMedia("(prefers-color-scheme: dark)").matches
  ) {
    document.body.classList.add("dark");
  }

  async function list(): Promise<void> {
    replays = await List();
  }

  async function main() {
    await list();

    // Listen for new replays
    // EventsOn('newReplay', async () => {
    //   await list();
    // });
  }

  async function play(fileName: string): Promise<void> {
    await Play(fileName);
  }

  async function playLocal(fileName: string): Promise<void> {
    await PlayLocal(fileName);
  }

  async function remove(fileName: string): Promise<void> {
    await Remove(fileName);
    removeTarget = "";
  }

  async function removeLocal(fileName: string): Promise<void> {
    await RemoveLocal(fileName);
    localRemoveTarget = "";
  }

  function formatLength(ticks) {
    const seconds = Math.floor(ticks / 8);
    const minutes = Math.floor(seconds / 60);
    const remainingSeconds = seconds % 60;
    const formattedMinutes = minutes < 10 ? `0${minutes}` : minutes;
    const formattedSeconds =
      remainingSeconds < 10 ? `0${remainingSeconds}` : remainingSeconds;
    return `${formattedMinutes}:${formattedSeconds}`;
  }

  function formatTime(timeString) {
    // Best guest regexes to convert the time string to a format that can be parsed by Date
    // Examples:
    // 13.04.2023 20:08 DD.MM.YYYY HH:MM
    // 2023/4/9上午 12:45 YYYY/MM/DDAM HH:MM
    // 08/04/2023 19:28 DD/MM/YYYY HH:MM
    // 4/12/2023 6:47 PM MM/DD/YYYY HH:MM AM/PM
    // 2023-03-13 오후 8:52 YYYY-MM-DD PM HH:MM

    const regex = /(\d{2})\.(\d{2})\.(\d{4}) (\d{2}):(\d{2})/;
    const regex2 = /(\d{4})\/(\d{1,2})\/(\d{1,2})/;
    const regex3 = /(\d{1,2})\/(\d{1,2})\/(\d{4}) (\d{2}):(\d{2})/;
    const regex4 = /(\d{1,2})\/(\d{1,2})\/(\d{4}) (\d{1,2}):(\d{2}) (AM|PM)/;
    const regex5 = /(\d{4})-(\d{1,2})-(\d{1,2})/;

    const guess1 = regex.exec(timeString);
    const guess2 = regex2.exec(timeString);
    const guess3 = regex3.exec(timeString);
    const guess4 = regex4.exec(timeString);
    const guess5 = regex5.exec(timeString);

    const originaltimeString = timeString;

    if (guess1) {
      timeString = `${pad(guess1[3])}-${pad(guess1[2])}-${pad(guess1[1])}T${pad(
        guess1[4]
      )}:${pad(guess1[5])}:00Z`;
    } else if (guess2) {
      timeString = `${pad(guess2[1])}-${pad(guess2[2])}-${pad(
        guess2[3]
      )}T00:00:00Z`;
    } else if (guess3) {
      if (Number(guess3[1]) > 12) {
        timeString = `${pad(guess3[3])}-${pad(guess3[2])}-${pad(
          guess3[1]
        )}T${pad(guess3[4])}:${pad(guess3[5])}:00Z`;
      } else {
        timeString = `${pad(guess3[3])}-${pad(guess3[1])}-${pad(
          guess3[2]
        )}T${pad(guess3[4])}:${pad(guess3[5])}:00Z`;
      }
    } else if (guess4) {
      if (Number(guess4[1]) > 12) {
        timeString = `${pad(guess4[3])}-${pad(guess4[2])}-${pad(
          guess4[1]
        )}T${pad(guess4[4])}:${pad(guess4[5])}:00Z`;
      } else {
        timeString = `${pad(guess4[3])}-${pad(guess4[1])}-${pad(
          guess4[2]
        )}T${pad(guess4[4])}:${pad(guess4[5])}:00Z`;
      }
    } else if (guess5) {
      timeString = `${pad(guess5[1])}-${pad(guess5[2])}-${pad(
        guess5[3]
      )}T00:00:00Z`;
    }

    const date = new Date(timeString);

    // Print the regexes that failed to help debug new formats
    if (isNaN(date.getTime())) {
      console.log("Invalid date: " + timeString);
      console.log("Original string: " + originaltimeString);

      if (!guess1) {
        console.log("Regex 1 failed");
      }
      if (!guess2) {
        console.log("Regex 2 failed");
      }
      if (!guess3) {
        console.log("Regex 3 failed");
      }
      if (!guess4) {
        console.log("Regex 4 failed");
      }
    }

    // Return in format: Apr 17, 2023
    const dateString = date.toLocaleDateString("en-GB", {
      month: "short",
      day: "numeric",
      year: "numeric",
    });
    return dateString;
  }

  function pad(num, size = 2) {
    let s = num + "";
    while (s.length < size) s = "0" + s;
    return s;
  }

  function formatMap(filename: string, players) {
    if (!filename.includes("\\")) {
      return filename;
    }
    // Example filename: "data:scenarios\multiplayer\twin_beach_2p_mkii\twin_beach_2p_mkii"
    // Another: "data:scenarios\multiplayer\(2) crossroads\(2) crossroads"
    // We want to return the last part of the path
    const parts = filename.split("\\");
    const mapName = parts[parts.length - 1];

    if (mapDetailsMap[mapName]) {
      return `
        <img src="${mapDetailsMap[mapName].url}" alt="${mapDetailsMap[mapName].name}" title="(${players.length}) ${mapDetailsMap[mapName].name}" style="width: 50px; height: auto;"/>
      `;
    } else {
      return mapName;
    }
  }

  const mapDetailsMap = {
    twin_beach_2p_mkii: {
      name: "Twin Beaches",
      url: twinBeachesIcon,
    },
    desert_village_2p_mkiii: {
      name: "Road to Tunis",
      url: roadToTunisIcon,
    },
    cliff_crossing_2p: {
      name: "Taranto Coastline",
      url: tarantoCoastlineIcon,
    },
    rails_and_sand_4p: {
      name: "Campbell's Convoy",
      url: campbellsConvoyIcon,
    },
    rural_town_4p: {
      name: "Pachino Farmlands",
      url: pachinoFarmlandsIcon,
    },
    torrente_4p_mkiii: {
      name: "Torrente",
      url: torrenteIcon,
    },
    rural_castle_4p: {
      name: "Aere Perennius",
      url: aerePerenniusIcon,
    },
    desert_airfield_6p_mkii: {
      name: "Gazala Landing Ground",
      url: gazalaLandingGroundIcon,
    },
    industrial_railyard_6p_mkii: {
      name: "L'Aquila",
      url: laquilaIcon,
    },
    winter_line_8p_mkii: {
      name: "Winter Line",
      url: winterLineIcon,
    },
    mountain_ruins_8p_mkii: {
      name: "Mignano Gap",
      url: mignanoGap8pIcon,
    },
    mountain_ruins_6p: {
      name: "Mignano Gap",
      url: mignanoGap6pIcon,
    },
    gardens_2p_mm: {
      name: "Gardens",
      url: gardensIcon,
    },
    "(2) crossroads": {
      name: "Crossroads",
      url: crossroadsIcon,
    },
    rural_town_2p_mkii: {
      name: "Pachino Stalemate",
      url: pachinoFarmlandsMkiiIcon,
    },
  };

  const factionIconMap = {
    AfrikaKorps: dakIcon,
    British: britishIcon,
    Wehrmacht: wehrmachtIcon,
    Americans: americanIcon,
  };

  function formatPlayers(players) {
    if (players.length === 0) {
      return "";
    }

    // Group players by team
    const teams = players.reduce((acc, player) => {
      if (!acc[player.Team]) {
        acc[player.Team] = [];
      }
      acc[player.Team].push(player);
      return acc;
    }, {});

    // For each team, create a string of faction icons
    let team1 = teams.First.map((player) => {
      return `<a href="https://steamcommunity.com/profiles/${
        player.SteamID
      }" target="_blank"><img class="faction-icon" src="${
        factionIconMap[player.Faction]
      }" title="${player.Name}"/></a>`;
    });

    let team2 = teams.Second.map((player) => {
      return `<a href="https://steamcommunity.com/profiles/${
        player.SteamID
      }" target="_blank"><img class="faction-icon" src="${
        factionIconMap[player.Faction]
      }" title="${player.Name}"/></a>`;
    });

    if (players.length === 2) {
      team1 = teams.First.map((player) => {
        return `<a href="https://steamcommunity.com/profiles/${
          player.SteamID
        }" target="_blank"><img class="faction-icon" src="${
          factionIconMap[player.Faction]
        }" title="${player.Name}"/></a> ${player.Name.trim()}`;
      });

      team2 = teams.Second.map((player) => {
        return `<a href="https://steamcommunity.com/profiles/${
          player.SteamID
        }" target="_blank"><img class="faction-icon" src="${
          factionIconMap[player.Faction]
        }" title="${player.Name}"/></a> ${player.Name.trim()}`;
      });
      return `${team1.join(" ")} <br> ${team2.join(" ")}`;
    }

    return `${team1.join(" ")} vs. ${team2.join(" ")}`;
  }

  main();
</script>

<main>
  <h1>Replay Manager</h1>
  <!-- Display table of replays -->
  {#if replays && replays.length > 0}
    <table>
      <thead>
        <tr>
          <!-- <th>Filename</th> -->
          <th>Map</th>
          <th>Players</th>
          <th>Game version</th>
          <th>Length</th>
          <th>Date</th>
          <th>Action</th>
        </tr>
      </thead>
      <tbody>
        {#each replays as replay}
          <tr>
            <!-- <td>{replay.Filename}</td> -->
            <td>{@html formatMap(replay.Map.Filename, replay.Players)}</td>
            <td>{@html formatPlayers(replay.Players)}</td>
            <td>{replay.Version}</td>
            <td>{formatLength(replay.Length)}</td>
            <td>{formatTime(replay.Timestamp)}</td>
            <td>
              <Row>
                <Col
                  ><Button on:click={() => play(replay.Filename)}
                    ><Icon src={mdiPlay} /></Button
                  ></Col
                >
                <Col
                  ><Button
                    on:click={() => {
                      modal_show();
                      removeTarget = replay.Filename;
                    }}
                    error><Icon src={mdiDelete} /></Button
                  ></Col
                >
              </Row></td
            >
          </tr>
        {/each}
      </tbody>
    </table>
  {:else}
    <p>No replays found</p>
  {/if}

  <!-- Display table of local replays
  {#if localReplays && Object.keys(localReplays).length > 0}
    <table>
      <thead>
        <tr>
          <th>Replay</th>
          <th>Date</th>
          <th>Game version</th>
        </tr>
      </thead>
      <tbody>
        {#each Object.keys(localReplays) as key (key)}
          <tr>
            {#if localReplays[key].includes("campaign")}
              <td>Local Campaign</td>
            {:else}
              <td>Local Skirmish</td>
            {/if}
            {#if localReplays[key].includes("campaign")}
              <td
                >{new Date(
                  parseInt(
                    localReplays[key].split("-")[3].replace(".rec", "")
                  ) / 1000000
                ).toLocaleString()}</td
              >
            {:else}
              <td
                >{new Date(
                  parseInt(
                    localReplays[key].split("-")[2].replace(".rec", "")
                  ) / 1000000
                ).toLocaleString()}</td
              >
            {/if}
            <td>N/A</td>
            <td>
              <Row>
                <Col
                  ><Button on:click={() => playLocal(localReplays[key])}
                    ><Icon src={mdiPlay} /></Button
                  ></Col
                >
                <Col
                  ><Button
                    on:click={() => {
                      modal_show_local();
                      localRemoveTarget = localReplays[key];
                    }}
                    error><Icon src={mdiDelete} /></Button
                  ></Col
                >
              </Row></td
            >
          </tr>
        {/each}
      </tbody>
    </table>
  {:else}
    <p>No local replays found</p>
  {/if} -->

  <Modal bind:open={openLocalRemove}>
    <Card>
      <h4 slot="header">Remove local replay</h4>

      <p>Are you sure you want to remove this replay?</p>

      <div slot="footer" class="is-right">
        <Button clear on:click={modal_hide_local}>Cancel</Button>
        <Button error on:click={() => removeLocal(localRemoveTarget)}
          >Remove</Button
        >
      </div>
    </Card>
  </Modal>

  <Modal bind:open={openRemove}>
    <Card>
      <h4 slot="header">Remove replay</h4>

      <p>Are you sure you want to remove this replay?</p>

      <div slot="footer" class="is-right">
        <Button clear on:click={modal_hide}>Cancel</Button>
        <Button error on:click={() => remove(removeTarget)}>Remove</Button>
      </div>
    </Card>
  </Modal>

  <!-- <Modal isOpen={openRemove} toggle={toggleRemove} {size}>
    <ModalHeader toggle={toggleRemove}>Remove replay</ModalHeader>
    <ModalBody>
      <p>Are you sure you want to remove this replay?</p>

      <p>It will be permanently deleted.</p>
    </ModalBody>
    <ModalFooter>
      <Button color="danger" on:click={() => remove(removeTarget)}
        >Remove</Button
      >
      <Button color="secondary" on:click={toggleRemove}>Cancel</Button>
    </ModalFooter>
  </Modal> -->
</main>

<style>
</style>
