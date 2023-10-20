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
  {key: 'state', value: 'State'},
  {key: 'name', value: 'Name'},
  {key: 'remote', value: 'Remote'},
  {key: 'actions', empty: true, width: '60px'},
]
var servers = []


onMount(() => refresh($query))
user.subscribe(() => refresh($query))

/**
 * @param query {URLSearchParams}
 */
function refresh(query) {
  // TODO if admin and query.all is false, query userServers, to show owned servers
  // TODO if user, show owned servers
  // TODO if admin, show all servers

  if (!$user) return

  /** @type {any} */
  var options = {}
  if ($user.isAdmin && !query.has('all')) options = {filter: `@collection.userServers.user ?= "${$user.id}" && @collection.userServers.server ?= id`}

  pb.collection('servers').getFullList(200, {...options, fields: 'id,name,username,host,port,state,created,updated'})
    .then(items => servers = items)
}

</script>

<Content>
  <DataTable zebra title={$query.has('all') ? 'Servers' : 'My Servers'} {headers} rows={servers}>
    <svelte:fragment slot="cell" let:row let:cell>
      {#if cell.key === "name" && $user.isAdmin}
        <Link href={'/servers/' + row.id}>{cell.value}</Link>
      {:else if cell.key === "remote"}
        <CopyButton text={sshTarget(row)} /> <Link href="ssh://{sshTarget(row)}">{sshTarget(row)}</Link>
      {:else if cell.key === "state"}
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
