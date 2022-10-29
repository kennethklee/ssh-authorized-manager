<script>
import {
  Content,
  DataTable,
  DataTableSkeleton,
  Toolbar,
  ToolbarContent,
  ToolbarSearch,
  OverflowMenu,
  OverflowMenuItem,
  Link,
  Button,
} from 'carbon-components-svelte'
import pb from '$lib/pocketbase.js'
import {formatDateTime} from '$lib/utils.js'


const headers = [
  {key: 'name', value: 'Name'},
  {key: 'updated', value: 'Last Updated'},
  {key: 'actions', empty: true, width: '60px'},
]
document.title = 'Public Keys | SSH Authorized Manager'
</script>

<Content>
  {#await pb.records.getFullList('publicKeys', 200)}
    <DataTableSkeleton />
  {:then items}
    <DataTable zebra title="My Public Keys" {headers} rows={items}>
      <svelte:fragment slot="cell" let:row let:cell>
        {#if cell.key === "name"}
          <Link href={'/publicKeys/' + row.id}>{cell.value || (row.publicKey.slice(0, 10) + '...')}</Link>
        {:else if cell.key === "updated"}
          <time datetime={formatDateTime(cell.value)}>{new Date(formatDateTime(cell.value)).toLocaleString()}</time>
        {:else if cell.key === "actions"}
          <OverflowMenu flipped>
            <OverflowMenuItem href={'/publicKeys/' + row.id}>Edit</OverflowMenuItem>
            <OverflowMenuItem on:click={() => {}}>Delete</OverflowMenuItem>
          </OverflowMenu>
        {:else}
          {cell.value}
        {/if}
      </svelte:fragment>

      <Toolbar>
        <ToolbarContent>
          <ToolbarSearch shouldFilterRows />
          <Button href="/publicKeys/new">Create public key</Button>
        </ToolbarContent>
      </Toolbar>
    </DataTable>
  {:catch error}
    <div>Error: {error}</div>
  {/await}
</Content>
