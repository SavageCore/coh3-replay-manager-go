<script context="module">
  export const mapDetailsMap = {
    rural_castle_4p: {
      name: 'Aere Perennius',
      url: aerePerenniusIcon,
    },
    rails_and_sand_4p: {
      name: "Campbell's Convoy",
      url: campbellsConvoyIcon,
    },
    '(2) crossroads': {
      name: 'Crossroads',
      url: crossroadsIcon,
    },
    gardens_2p_mm: {
      name: 'Gardens',
      url: gardensIcon,
    },
    desert_airfield_6p_mkii: {
      name: 'Gazala Landing Ground',
      url: gazalaLandingGroundIcon,
    },
    industrial_railyard_6p_mkii: {
      name: "L'Aquila",
      url: laquilaIcon,
    },
    mountain_ruins_6p: {
      name: 'Mignano Gap 3v3',
      url: mignanoGap6pIcon,
    },
    mountain_ruins_8p_mkii: {
      name: 'Mignano Gap 4v4',
      url: mignanoGap8pIcon,
    },
    rural_town_4p: {
      name: 'Pachino Farmlands',
      url: pachinoFarmlandsIcon,
    },
    rural_town_2p_mkii: {
      name: 'Pachino Stalemate',
      url: pachinoFarmlandsMkiiIcon,
    },
    desert_village_2p_mkiii: {
      name: 'Road to Tunis',
      url: roadToTunisIcon,
    },
    cliff_crossing_2p: {
      name: 'Taranto Coastline',
      url: tarantoCoastlineIcon,
    },
    torrente_4p_mkiii: {
      name: 'Torrente',
      url: torrenteIcon,
    },
    twin_beach_2p_mkii: {
      name: 'Twin Beaches',
      url: twinBeachesIcon,
    },
    winter_line_8p_mkii: {
      name: 'Winter Line',
      url: winterLineIcon,
    },
  };
</script>

<script lang="ts">
  import { GetGameVersion, List } from '../wailsjs/go/main/App.js';
  import { Card, Container, Field } from 'svelte-chota';
  import { CollapsibleCard } from 'svelte-collapsible';
  import 'chota';

  import Table from './Table.svelte';
  import { replays as replayStore } from './data-store.js';

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

  let replays;
  let unfilteredReplays;
  let currentGameVersion;
  let open = false;

  const playerCounts = [2, 4, 6, 8];

  if (
    window.matchMedia &&
    window.matchMedia('(prefers-color-scheme: dark)').matches
  ) {
    document.body.classList.add('dark');
  }

  async function list(): Promise<void> {
    replays = await List();
    unfilteredReplays = replays;

    replayStore.set(replays);
  }

  async function main() {
    currentGameVersion = await GetGameVersion();

    await list();
  }

  async function filter(property, needle) {
    replays = unfilteredReplays.filter((replay) => {
      if (needle === 'all') {
        return true;
      }

      if (property.includes('.')) {
        const [parent, child] = property.split('.');

        return replay[parent][child] === needle;
      }

      if (property === 'Players') {
        return replay.Players.length === needle;
      }

      return replay[property] === needle;
    });

    replayStore.set(replays);
  }

  main();
</script>

<main>
  <Container>
    <CollapsibleCard
      {open}
      on:open={() => {
        const header = document.querySelector('#filter-header');
        const h2 = document.createElement('h2');
        h2.innerHTML = header.innerHTML;
        h2.id = header.id;
        header.replaceWith(h2);
      }}
      on:close={() => {
        const header = document.querySelector('#filter-header');
        const p = document.createElement('p');
        p.innerHTML = header.innerHTML;
        p.id = header.id;
        header.replaceWith(p);
      }}
    >
      <p slot="header" id="filter-header">Filters</p>
      <div slot="body">
        <Field label="Map">
          <select
            on:change={(event) => {
              filter('Map.Filename', event.target.value);
            }}
          >
            <option value="all" selected>All</option>
            {#each Object.entries(mapDetailsMap) as [key, map]}
              <option value={key}>{map.name}</option>
            {/each}
          </select>
        </Field>
        <Field label="Player count">
          <select
            on:change={(event) => {
              if (event.target.value === 'all') {
                filter('Players', 'all');

                return;
              }

              filter('Players', Number(event.target.value));
            }}
          >
            <option value="all" selected>All</option>
            {#each playerCounts as count}
              <option value={count}>{count}</option>
            {/each}
          </select>
        </Field>
        <Field label="Hide incompatible">
          <input
            type="checkbox"
            on:change={(event) => {
              // If unchecked, show all replays
              if (!event.target.checked) {
                filter('Version', 'all');

                return;
              }
              // If checked, filter out replays that aren't compatible with the current game version
              filter('Version', Number(currentGameVersion));
            }}
          />
        </Field>
      </div>
    </CollapsibleCard>

    {#if replays && replays.length > 0}
      <Table tableData={$replayStore} />
    {/if}
  </Container>
</main>

<style>
</style>
