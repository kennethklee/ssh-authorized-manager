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
  pb.collection('servers').getFullList(200, {fields: 'id,name,username,host,port,state,created,updated'})
    .then(items => {
      allServers = items

      pb.collection('userServers').getFullList(200, {filter: `user="${user.id}"`})
        .then(items => {
          userServers = items
          selectedRowIds = items.map(us => us.server)
        })
    })
}

function handleSave(ev) {
  // diff original and selected server ids
  const addedServerIds = selectedRowIds.filter(id => !userServers.some(us => us.server === id))
  const removedUserServers = userServers.filter(us => !selectedRowIds.includes(us.server))

  // add new user servers
  addedServerIds.forEach(id => {
    pb.collection('userServers').create({user: user.id, server: id}, {$autoCancel: false})
  })
  // remove old user servers
  removedUserServers.forEach(us => {
    pb.collection('userServers').delete(us.id, {$autoCancel: false})
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

