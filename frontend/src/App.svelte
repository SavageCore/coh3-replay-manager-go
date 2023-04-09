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
    const listResponse = await List();
    replays = listResponse.Replays;
    localReplays = listResponse.LocalReplays;
  }

  async function main() {
    await list();

    // Listen for new replays
    // EventsOn('newReplay', async () => {
    //   await list();
    // });
  }

  async function play(id: string): Promise<void> {
    await Play(id);
  }

  async function playLocal(id: string): Promise<void> {
    await PlayLocal(id);
  }

  async function remove(id: string): Promise<void> {
    console.log(id);
    await Remove(id);
    removeTarget = "";
  }

  async function removeLocal(id: string): Promise<void> {
    await RemoveLocal(id);
    localRemoveTarget = "";
  }

  main();
</script>

<main>
  <h1>Replay Manager</h1>
  <!-- <p>openLocalRemove={openLocalRemove}</p> -->
  <!-- Display table of replays -->
  {#if replays && Object.keys(replays).length > 0}
    <table>
      <thead>
        <tr>
          <th>Replay</th>
          <th>Game version</th>
          <th>Action</th>
        </tr>
      </thead>
      <tbody>
        {#each Object.keys(replays) as key (key)}
          <tr>
            <td>{replays[key].split("|")[0]}</td>
            <td>{replays[key].split("|")[1]}</td>
            <td>
              <Row>
                <Col
                  ><Button on:click={() => play(key)}
                    ><Icon src={mdiPlay} /></Button
                  ></Col
                >
                <Col
                  ><Button
                    on:click={() => {
                      modal_show();
                      removeTarget = key;
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

  <!-- Display table of local replays -->
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
  {/if}

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
