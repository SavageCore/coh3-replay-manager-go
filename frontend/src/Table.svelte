<script>
  // Import map icons
  import aerePerenniusIcon from './assets/icons/maps/aere_perennius.webp';
  import campbellsConvoyIcon from './assets/icons/maps/campbells_convoy.webp';
  import gardensIcon from './assets/icons/maps/gardens.png';
  import crossroadsIcon from './assets/icons/maps/gardens.png';
  import gazalaLandingGroundIcon from './assets/icons/maps/gazala_landing_ground.webp';
  import laquilaIcon from './assets/icons/maps/laquila.webp';
  import mignanoGap6pIcon from './assets/icons/maps/mignano_gap_6p.png';
  import mignanoGap8pIcon from './assets/icons/maps/mignano_gap.webp';
  import pachinoFarmlandsIcon from './assets/icons/maps/pachino_farmlands.webp';
  import pachinoFarmlandsMkiiIcon from './assets/icons/maps/pachino_farmlands_mkii.png';
  import roadToTunisIcon from './assets/icons/maps/road_to_tunis.webp';
  import tarantoCoastlineIcon from './assets/icons/maps/taranto_coastline.webp';
  import torrenteIcon from './assets/icons/maps/torrente.webp';
  import twinBeachesIcon from './assets/icons/maps/twin_beaches.webp';
  import winterLineIcon from './assets/icons/maps/winter_line.webp';

  // Import faction icons
  import dakIcon from './assets/icons/factions/dak.webp';
  import britishIcon from './assets/icons/factions/british.webp';
  import wehrmachtIcon from './assets/icons/factions/german.webp';
  import americanIcon from './assets/icons/factions/american.webp';

  import { Button, Row, Col, Icon } from 'svelte-chota';
  import { mdiDelete, mdiPlay } from '@mdi/js';
  import { Play } from '../wailsjs/go/main/App.js';

  export let tableData;

  let sorted = false;

  // Function to sort a table column
  //   event: the event object
  //   column: the column to sort by
  //   Flip the sort order if the column is already sorted
  const sort = (event, column) => {
    if (sorted) {
      tableData.sort((a, b) => {
        let aSort = a[column];
        let bSort = b[column];

        if (column.includes('.')) {
          const parts = column.split('.');
          aSort = a[parts[0]][parts[1]];
          bSort = b[parts[0]][parts[1]];
        }

        if (column === 'Players') {
          aSort = a[column].length;
          bSort = b[column].length;
        }

        if (aSort < bSort) {
          return 1;
        }
        if (aSort > bSort) {
          return -1;
        }
        return 0;
      });
      sorted = false;
    } else {
      tableData.sort((a, b) => {
        let aSort = a[column];
        let bSort = b[column];

        if (column.includes('.')) {
          const parts = column.split('.');
          aSort = a[parts[0]][parts[1]];
          bSort = b[parts[0]][parts[1]];
        }

        if (column === 'Players') {
          aSort = a[column].length;
          bSort = b[column].length;
        }

        if (aSort < bSort) {
          return -1;
        }
        if (aSort > bSort) {
          return 1;
        }
        return 0;
      });
      sorted = true;
    }

    tableData = tableData;
  };

  async function play(fileName) {
    await Play(fileName);
  }

  const formatLength = (ticks) => {
    const seconds = Math.floor(ticks / 8);
    const minutes = Math.floor(seconds / 60);
    const remainingSeconds = seconds % 60;
    const formattedMinutes = minutes < 10 ? `0${minutes}` : minutes;
    const formattedSeconds =
      remainingSeconds < 10 ? `0${remainingSeconds}` : remainingSeconds;
    return `${formattedMinutes}:${formattedSeconds}`;
  };

  const displayMap = (filename, players) => {
    if (mapDetailsMap[filename]) {
      return `
        <img src="${mapDetailsMap[filename].url}" alt="${mapDetailsMap[filename].name}" title="(${players.length}) ${mapDetailsMap[filename].name}" style="width: 50px; height: auto;"/>
      `;
    } else {
      return filename;
    }
  };

  const formatPlayers = (players) => {
    if (players.length === 0) {
      return '';
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
      return `${team1.join(' ')} <br> ${team2.join(' ')}`;
    }

    return `${team1.join(' ')} vs. ${team2.join(' ')}`;
  };

  const mapDetailsMap = {
    twin_beach_2p_mkii: {
      name: 'Twin Beaches',
      url: twinBeachesIcon,
    },
    desert_village_2p_mkiii: {
      name: 'Road to Tunis',
      url: roadToTunisIcon,
    },
    cliff_crossing_2p: {
      name: 'Taranto Coastline',
      url: tarantoCoastlineIcon,
    },
    rails_and_sand_4p: {
      name: "Campbell's Convoy",
      url: campbellsConvoyIcon,
    },
    rural_town_4p: {
      name: 'Pachino Farmlands',
      url: pachinoFarmlandsIcon,
    },
    torrente_4p_mkiii: {
      name: 'Torrente',
      url: torrenteIcon,
    },
    rural_castle_4p: {
      name: 'Aere Perennius',
      url: aerePerenniusIcon,
    },
    desert_airfield_6p_mkii: {
      name: 'Gazala Landing Ground',
      url: gazalaLandingGroundIcon,
    },
    industrial_railyard_6p_mkii: {
      name: "L'Aquila",
      url: laquilaIcon,
    },
    winter_line_8p_mkii: {
      name: 'Winter Line',
      url: winterLineIcon,
    },
    mountain_ruins_8p_mkii: {
      name: 'Mignano Gap',
      url: mignanoGap8pIcon,
    },
    mountain_ruins_6p: {
      name: 'Mignano Gap',
      url: mignanoGap6pIcon,
    },
    gardens_2p_mm: {
      name: 'Gardens',
      url: gardensIcon,
    },
    '(2) crossroads': {
      name: 'Crossroads',
      url: crossroadsIcon,
    },
    rural_town_2p_mkii: {
      name: 'Pachino Stalemate',
      url: pachinoFarmlandsMkiiIcon,
    },
  };

  const factionIconMap = {
    AfrikaKorps: dakIcon,
    British: britishIcon,
    Wehrmacht: wehrmachtIcon,
    Americans: americanIcon,
  };
</script>

<table>
  <thead>
    <tr>
      <th
        on:click={() => {
          sort(event, 'Map.Filename');
        }}>Map</th
      >
      <th
        on:click={() => {
          sort(event, 'Players');
        }}>Players</th
      >
      <th
        on:click={() => {
          sort(event, 'Version');
        }}>Game version</th
      >
      <th
        on:click={() => {
          sort(event, 'Length');
        }}>Length</th
      >
      <th
        on:click={() => {
          sort(event, 'Timestamp');
        }}>Date</th
      >
      <th>Action</th>
    </tr>
  </thead>
  <tbody>
    {#each tableData as replay}
      <tr>
        <td>{@html displayMap(replay.Map.Filename, replay.Players)}</td>
        <td>{@html formatPlayers(replay.Players)}</td>
        <td>{replay.Version}</td>
        <td>{formatLength(replay.Length)}</td>
        <td>{replay.Timestamp}</td>
        <td>
          <Row>
            <Col
              ><Button on:click={() => play(replay.Filename)}
                ><Icon src={mdiPlay} /></Button
              ></Col
            >
            <Col><Button error><Icon src={mdiDelete} /></Button></Col>
          </Row></td
        >
      </tr>
    {/each}
  </tbody>
</table>

<style>
  table,
  th,
  td {
    border: 1px solid #ccc;
    border-collapse: collapse;
    margin-bottom: 10px;
  }
</style>
