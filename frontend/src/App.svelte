<script lang="ts">
  import { List } from '../wailsjs/go/main/App.js';
  import { Container } from 'svelte-chota';
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
  <!-- Show a table of replays -->
  <Container>
    <h1>Replays</h1>

    {#if replays && replays.length > 0}
      <Table tableData={$replayStore} />
    {/if}
  </Container>
</main>

<style>
</style>
