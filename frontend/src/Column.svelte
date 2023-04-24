<script>
  export let column;
  export let sortColumn;
  export let sortDirection;
  export let sort;
</script>

<th
  class:sorted={sortColumn === column.field}
  class:asc={sortDirection === 'asc'}
  class:desc={sortDirection === 'desc'}
  on:click={() => {
    [sortColumn, sortDirection] = sort(column.field);
  }}>{column.label || column.field}</th
>

<style>
  .sorted {
    background-color: var(--bg-secondary-color);
  }

  .sorted.asc:after {
    content: ' ▲';
  }

  .sorted.desc:after {
    content: ' ▼';
  }
</style>
