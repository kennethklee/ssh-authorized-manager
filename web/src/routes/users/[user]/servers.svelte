<script>
import {
  Content,
  DataTable,
  Toolbar,
  ToolbarContent,
  ToolbarSearch,
  Link,
  Button,
  InlineNotification
} from 'carbon-components-svelte'
import {onMount} from 'svelte'

import pb from '$lib/pocketbase.js'


export var user = {}
const headers = [
  {key: 'name', value: 'Name'},
  {key: 'remote', value: 'Remote'},
]
var allServers = []
var userServers = []
var selectedRowIds = []
let saved = false



onMount(() => refresh())

function refresh(query) {
  /** @type {any} */
  pb.records.getFullList('servers', 200)
    .then(items => {
      allServers = items

      pb.records.getFullList('userServers', 200, {filter: `userId="${user.id}"`})
        .then(items => {
          userServers = items
          selectedRowIds = items.map(us => us.serverId)
        })
    })
}

function handleSave(ev) {
  // diff original and selected server ids
  const addedServerIds = selectedRowIds.filter(id => !userServers.some(us => us.serverId === id))
  const removedUserServers = userServers.filter(us => !selectedRowIds.includes(us.serverId))

  // add new user servers
  addedServerIds.forEach(id => {
    pb.records.create('userServers', {userId: user.id, serverId: id}, {$autoCancel: false})
  })
  // remove old user servers
  removedUserServers.forEach(us => {
    pb.records.delete('userServers', us.id, {$autoCancel: false})
  })

  saved = true
  setTimeout(() => saved = false, 1500)
}

// TODO add special options to particular servers
</script>

<style>
  nav {
    display: flex;
    justify-content: space-between;
  }
</style>

<Content>
  <nav>
    <header>
      <h2>{user.profile?.name || user.email} servers</h2>
    </header>
    <article>
      <Button on:click={handleSave}>Save</Button>
    </article>
  </nav>

  {#if saved}
    <InlineNotification kind="success" title="Saved" subtitle="Your changes have been saved." hideCloseButton />
  {/if}

  <DataTable batchSelection bind:selectedRowIds sortable zebra title="Server Access" {headers} rows={allServers}>
    <svelte:fragment slot="cell" let:row let:cell>
      {#if cell.key === "name"}
        <Link href={'/servers/' + row.id}>{cell.value}</Link>
      {:else if cell.key === "remote"}
        {#if row.port === 22}
          {row.username}@{row.host}
        {:else}
          {row.username}@{row.host}:{row.port}
        {/if}
      {:else}
        {cell.value}
      {/if}
    </svelte:fragment>

    <Toolbar>
      <ToolbarContent>
        <ToolbarSearch shouldFilterRows />
      </ToolbarContent>
    </Toolbar>
  </DataTable>
</Content>

