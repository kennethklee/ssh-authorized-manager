<script>
import {
  Content,
  DataTable,
  Toolbar,
  ToolbarContent,
  ToolbarSearch,
  OverflowMenu,
  OverflowMenuItem,
  Link,
} from 'carbon-components-svelte'
import {onMount} from 'svelte'

import pb from '$lib/pocketbase.js'
import {formatDateTime} from '$lib/utils.js'


const headers = [
  {key: 'name', value: 'Name'},
  {key: 'email', value: 'Email'},
  {key: 'updated', value: 'Last Updated'},
  {key: 'actions', empty: true, width: '60px'},
]
var users = []


onMount(() => refresh())

async function refresh() {
  users = await pb.collection('users').getFullList(200, {$autoCancel: false})
}
</script>

<Content>
  <DataTable zebra title="Users" {headers} rows={users}>
    <svelte:fragment slot="cell" let:row let:cell>
      {#if cell.key === "updated"}
        <time datetime={formatDateTime(cell.value)}>{new Date(formatDateTime(cell.value)).toLocaleString()}</time>
      {:else if cell.key === "actions"}
        <OverflowMenu flipped>
          <OverflowMenuItem href={`/users/${row.id}/servers`}>Manage Servers</OverflowMenuItem>
        </OverflowMenu>
      {:else}
        {cell.value || 'N/A'}
      {/if}
    </svelte:fragment>

    <Toolbar>
      <ToolbarContent>
        <ToolbarSearch shouldFilterRows />
      </ToolbarContent>
    </Toolbar>
  </DataTable>
</Content>
