<script>
import {onMount, onDestroy} from 'svelte'
import {
  Modal,
  Content,
  DataTable,
  OverflowMenu,
  Toolbar,
  ToolbarContent,
  Button,
} from 'carbon-components-svelte'
import Repeat from 'carbon-icons-svelte/lib/Repeat.svelte'
import pb from '$lib/pocketbase.js'
import {formatDateTime} from '$lib/utils.js'


export var server = {port: 22, usePassword: true}
const headers = [
  {key: 'created', value: 'Date', width: '200px'},
  {key: 'type', value: 'Type', width: '100px'},
  {key: 'message', value: 'Message'},
  {key: 'actions', empty: true, width: '60px'},
]
var logs = []
var isDetailsOpen = false
var details = null


onMount(() => {
  refresh()

  pb.realtime.subscribe('serverLogs', function (e) {
    if (e.record.serverId !== server.id) return ;
    if (e.action === 'create' && e.record.serverId == server.id) {
      logs = [e.record, ...logs]
    } else if (e.action === 'delete' && e.record.serverId == server.id) {
      logs = logs.filter(l => l.id !== e.record.id)
    }
  })
})

onDestroy(() => {
  pb.realtime.unsubscribe('serverLogs')
})

function refresh() {
  pb.records.getList('serverLogs', 1, 200, {sort: '-created', filter: `serverId="${server.id}"`})
    .then(data => logs = data.items)
}

function showDetails(data) {
  details = data
  isDetailsOpen = true
}

function handleDetailsTransition(ev) {
  if (!isDetailsOpen) {
    details = null
  }
}

function deleteMessage(id) {
  isDetailsOpen = false
  pb.records.delete('serverLogs', id)
}

function deletePastMessages() {
  // delete all messages except the last one
  logs.slice(1).forEach(l => pb.records.delete('serverLogs', l.id))
}

function trustHostKey() {
  pb.records.update('servers', server.id, {hostKey: details.payload})
  deleteMessage(details.id)
}

function trustHostName() {
  pb.records.update('servers', server.id, {hostname: details.payload})
  deleteMessage(details.id)
}
</script>

<Modal
    bind:open={isDetailsOpen}
    modalHeading="Log Details"
    primaryButtonText="Close"
    on:click:button--primary={() => isDetailsOpen = false}
    on:transitionend={handleDetailsTransition}>
  <p>{new Date(details?.created+'z').toLocaleString()}</p>
  <p><strong>{details?.message}</strong></p>
  {#if details?.type === 'hostKey'}
  <Button kind="tertiary" on:click={trustHostKey}>Trust server fingerprint</Button>
  <Button kind="danger-ghost" on:click={() => deleteMessage(details.id)}>Not trusted</Button>
  {:else if details?.type === 'hostName'}
  <Button kind="tertiary" on:click={trustHostName}>Trust server hostname</Button>
  <Button kind="danger-ghost" on:click={() => deleteMessage(details.id)}>Not trusted</Button>
  {:else}
  <p>{details?.payload}</p>
  {/if}
</Modal>

<Content>
  <Button style="float: right;">Force Sync</Button>
  <Button href="/servers/{server.id}">Back to {server.name}</Button>

  <DataTable zebra title="Sync logs for {server.name}" {headers} rows={logs}>
    <svelte:fragment slot="cell" let:row let:cell>
      {#if cell.key === "message"}
        {cell.value}
      {:else if cell.key === "created"}
        <time datetime={formatDateTime(cell.value)}>{new Date(formatDateTime(cell.value)).toLocaleString()}</time>
      {:else if cell.key === "actions"}
        <OverflowMenu on:click={() => showDetails(row)} />
      {:else}
        {cell.value}
      {/if}
    </svelte:fragment>

    <Toolbar size="sm">
      <ToolbarContent>
        <Button kind="ghost" icon={Repeat} iconDescription="Refresh" on:click={() => refresh()} />
        <Button kind="danger-ghost" on:click={() => deletePastMessages()}>Clear Old Logs</Button>
      </ToolbarContent>
    </Toolbar>

  </DataTable>
</Content>
