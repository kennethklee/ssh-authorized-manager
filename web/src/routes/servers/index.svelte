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
  Button,
  CopyButton,
} from 'carbon-components-svelte'
import { onMount } from 'svelte'

import pb from '$lib/pocketbase.js'
import {user, query} from '$lib/stores.js'
import {sshTarget} from '$lib/utils'
import LogState from './_logState.svelte'


const headers = [
  {key: 'status', value: 'Last Sync'},
  {key: 'name', value: 'Name'},
  {key: 'remote', value: 'Remote'},
  {key: 'actions', empty: true, width: '60px'},
]
var servers = []


onMount(() => refresh($query))

/**
 * @param query {URLSearchParams}
 */
function refresh(query) {
  // TODO if admin and query.all is false, query userServers, to show owned servers
  // TODO if user, show owned servers
  // TODO if admin, show all servers

  /** @type {any} */
  var options = {}
  if ($user.isAdmin && !query.has('all')) options = {filter: `@collection.userServers.userId="${$user.id}" && @collection.userServers.serverId=id`}

  pb.collection('servers').getFullList(200, options)
    .then(items => {
      servers = items

      // TODO need a better way to group one log per server on server-side instead of client-side
      return Promise.all(items.map(s => fetchLastServerLog(s.id)))
    })
    .then(logs => {
      // make a map of serverId to messages
      const state = {}
      logs.forEach(log => state[log.serverId] = log)
      servers.forEach(server => server.lastLog = state[server.id])
      servers = servers // force datatable update
    })
}

function fetchLastServerLog(id) {
  return pb.collection('serverLogs').getList(1, 1, {filter: `serverId='${id}'`, sort: '-created', $autoCancel: false})
    .then(data => data.items[0])
}

</script>

<Content>
  <DataTable zebra title={$query.has('all') ? 'Servers' : 'My Servers'} {headers} rows={servers}>
    <svelte:fragment slot="cell" let:row let:cell>
      {#if cell.key === "name" && !$user.isUser}
        <Link href={'/servers/' + row.id}>{cell.value}</Link>
      {:else if cell.key === "remote"}
        <CopyButton text={sshTarget(row)} /> {sshTarget(row)}
      {:else if cell.key === "status"}
        <LogState {row} />
      {:else if cell.key === "actions"}
        {#if $user?.isAdmin}
        <OverflowMenu flipped>
          <OverflowMenuItem href="/servers/{row.id}/logs">View sync logs</OverflowMenuItem>
          <OverflowMenuItem href="/servers/{row.id}/users">User Access</OverflowMenuItem>
          <OverflowMenuItem href="/servers/{row.id}">Edit</OverflowMenuItem>
        </OverflowMenu>
        {/if}
      {:else}
        {cell.value}
      {/if}
    </svelte:fragment>

    <Toolbar>
      <ToolbarContent>
        <ToolbarSearch shouldFilterRows />
        {#if $user?.isAdmin}
          <Button href="/servers/new">Create server</Button>
        {/if}
      </ToolbarContent>
    </Toolbar>
  </DataTable>
</Content>
