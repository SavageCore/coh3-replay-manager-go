<script lang="ts">
  import { List } from '../wailsjs/go/main/App.js';
  import 'chota';

  import Table from './Table.svelte';
  import { replays as replayStore } from './data-store.js';

  let replays;

  if (
    window.matchMedia &&
    window.matchMedia('(prefers-color-scheme: dark)').matches
  ) {
    document.body.classList.add('dark');
  }

  async function list(): Promise<void> {
    replays = await List();

    replayStore.set(replays);
  }

  async function main() {
    await list();
  }

  main();
</script>

<main>
  <h1>Replay Manager</h1>
  {#if replays && replays.length > 0}
    <Table tableData={$replayStore} />
  {/if}
</main>

<style>
</style>
