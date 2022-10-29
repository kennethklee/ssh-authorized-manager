<script>
import {
  Content,
  DataTable,
  Toolbar,
  ToolbarContent,
  ToolbarSearch,
  ToolbarMenu,
  ToolbarMenuItem,
  Select,
  SelectItem,
  OverflowMenu,
  OverflowMenuItem,
  Link,
  Button,
} from 'carbon-components-svelte'


const headers = [
  {key: 'name', value: 'Name'},
  {key: 'email', value: 'Email'},
  {key: 'actions', empty: true, width: '60px'},
]

export var server = {port: 22, usePassword: true}
var users = []

</script>

<Content>
  <Button href="/servers/{server.id}">Back to {server.name}</Button>

  <DataTable batchSelection sortable zebra title="Users" {headers} rows={users}>
    <svelte:fragment slot="cell" let:row let:cell>
      {#if cell.key === "actions"}
        <OverflowMenu flipped>
          <OverflowMenuItem href={`/users/${row.id}/servers`}>Manage Servers</OverflowMenuItem>
          <OverflowMenuItem on:click={() => {}}>Delete</OverflowMenuItem>
        </OverflowMenu>
      {:else}
        {cell.value || 'N/A'}
      {/if}
    </svelte:fragment>

    <Toolbar size="sm">
      <ToolbarContent>
        <ToolbarSearch shouldFilterRows />
      </ToolbarContent>
    </Toolbar>
  </DataTable>
</Content>
