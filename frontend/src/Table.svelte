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
  import { mdiAlert, mdiDelete, mdiPlay } from '@mdi/js';

  import { Play, Remove, GetGameVersion } from '../wailsjs/go/main/App.js';
  import Column from './Column.svelte';

  export let tableData;

  let sorted = false;
  let sortColumn = 'Version';
  let sortDirection = 'desc';
  let currentGameVersion;
  async function main() {
    currentGameVersion = await GetGameVersion();

    // Sort the table by version by default
    sort('Version', 'desc');
  }

  main();

  // Function to sort tableData by a property with an optional direction
  const sort = (property, direction) => {
    if (!direction) {
      direction = sortDirection === 'asc' ? 'desc' : 'asc';
    }

    tableData.sort((a, b) => {
      let aSort = a[property];
      let bSort = b[property];

      if (property.includes('.')) {
        const parts = property.split('.');
        aSort = a[parts[0]][parts[1]];
        bSort = b[parts[0]][parts[1]];

        if (property === 'Map.Filename') {
          aSort = mapDetailsMap[a[parts[0]][parts[1]]].name;
          bSort = mapDetailsMap[b[parts[0]][parts[1]]].name;
        }
      }

      if (property === 'Players') {
        aSort = a[property].length;
        bSort = b[property].length;
      }

      if (aSort < bSort) {
        return direction === 'asc' ? -1 : 1;
      } else if (aSort > bSort) {
        return direction === 'asc' ? 1 : -1;
      } else {
        return 0;
      }
    });

    sorted = true;
    sortColumn = property;
    sortDirection = direction;

    tableData = tableData;

    return [sortColumn, sortDirection];
  };

  async function play(fileName) {
    await Play(fileName);
  }

  async function remove(fileName) {
    await Remove(fileName);
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
      <Column
        column={{ field: 'Map.Filename', label: 'Map' }}
        {sortColumn}
        {sortDirection}
        {sort}
      />
      <Column
        column={{ field: 'Players' }}
        {sortColumn}
        {sortDirection}
        {sort}
      />
      <Column
        column={{ field: 'Version' }}
        {sortColumn}
        {sortDirection}
        {sort}
      />
      <Column
        column={{ field: 'Length' }}
        {sortColumn}
        {sortDirection}
        {sort}
      />
      <Column
        column={{ field: 'Timestamp', label: 'Date' }}
        {sortColumn}
        {sortDirection}
        {sort}
      />
      <th>Action</th>
    </tr>
  </thead>
  <tbody>
    {#each tableData as replay}
      <tr>
        <td class="center"
          >{@html displayMap(replay.Map.Filename, replay.Players)}</td
        >
        <td>{@html formatPlayers(replay.Players)}</td>
        <td
          class:version-mismatch={currentGameVersion != replay.Version}
          title={currentGameVersion != replay.Version
            ? 'Version Mismatch - playback unlikely'
            : ''}
        >
          {#if currentGameVersion != replay.Version}
            <Icon src={mdiAlert} color="yellow" />
          {/if}
          {replay.Version}
        </td>
        <td>{formatLength(replay.Length)}</td>
        <td>{replay.Timestamp}</td>
        <td class="center">
          <Row>
            <Col
              ><Button
                on:click={() => play(replay.Filename)}
                icon={mdiPlay}
              /></Col
            >
            <Col
              ><Button
                error
                on:click={() => remove(replay.Filename)}
                icon={mdiDelete}
              /></Col
            >
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
    border: 1px solid var(--border-color);
    border-collapse: collapse;
    margin-bottom: 10px;
  }

  .version-mismatch {
    color: var(--color-error);
  }

  .version-mismatch:hover {
    color: var(--color-errorShade);
  }

  thead {
    background-color: var(--bg-secondary-color);
    color: var(--font-color);
    border-bottom: 3px solid var(--border-color);
  }

  table {
    background-color: #1b1a18;
  }

  tr:hover {
    background-color: #6f6f6f;
  }

  /* No hover style on thead */
  thead tr:hover {
    background-color: transparent;
  }

  td.center {
    text-align: center;
    padding-left: 0;
  }

  td,
  th {
    padding-left: 10px;
  }
</style>
