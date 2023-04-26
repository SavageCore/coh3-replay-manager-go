<script>
  export let showModal;
  export let closeClass = '';
  export let closeButtonText = 'Close';
  export let closeFunction = null;
  export let confirmClass = '';
  export let confirmButtonText = 'Confirm';
  export let confirmFunction = null;
  export let buttons = true;
  export let padding = null;

  import { Button, Row } from 'svelte-chota';

  let dialog;

  $: if (dialog && showModal) dialog.showModal();
</script>

<!-- svelte-ignore a11y-click-events-have-key-events -->
<dialog
  bind:this={dialog}
  on:close={() => (showModal = false)}
  on:click|self={() => dialog.close()}
>
  <div on:click|stopPropagation style:padding>
    <slot name="header" />
    <slot />
    <Row>
      <!-- svelte-ignore a11y-autofocus -->
      {#if buttons}
        <Button
          autofocus
          class={closeClass}
          on:click={() => closeFunction(dialog)}>{closeButtonText}</Button
        >
        {#if confirmFunction}
          <Button
            error
            class={confirmClass}
            on:click={() => confirmFunction(dialog)}>{confirmButtonText}</Button
          >
        {/if}
      {/if}
    </Row>
  </div>
</dialog>

<style>
  dialog {
    max-width: 32em;
    border-radius: 0.25em;
    border: none;
    padding: 0;
    background-color: var(--bg-color);
    color: var(--font-color);
  }

  dialog::backdrop {
    background: rgba(0, 0, 0, 0.3);
  }

  dialog > div {
    padding: 1em;
  }

  dialog[open] {
    animation: zoom 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
  }

  dialog[open]::backdrop {
    animation: fade 0.2s ease-out;
  }

  @keyframes zoom {
    from {
      transform: scale(0.95);
    }
    to {
      transform: scale(1);
    }
  }

  @keyframes fade {
    from {
      opacity: 0;
    }
    to {
      opacity: 1;
    }
  }
</style>
