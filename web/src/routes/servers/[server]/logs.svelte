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
  Tag,
} from 'carbon-components-svelte'
import * as timeago from 'timeago.js'
import Repeat from 'carbon-icons-svelte/lib/Repeat.svelte'
import pb from '$lib/pocketbase.js'
import {formatDateTime} from '$lib/utils.js'


export var server = {port: 22, usePassword: true}
const headers = [
  {key: 'ago', empty: true, width: '150px'},
  {key: 'created', value: 'Date', width: '200px'},
  {key: 'message', value: 'Message'},
  {key: 'actions', empty: true, width: '60px'},
]
var logs = []
var isDetailsOpen = false
var details = null


onMount(() => {
  refresh()

  pb.realtime.subscribe('serverLogs', function (e) {
    if (e.record.server !== server.id) return ;
    if (e.action === 'create' && e.record.server == server.id) {
      logs = [e.record, ...logs]
    } else if (e.action === 'delete' && e.record.server == server.id) {
      logs = logs.filter(l => l.id !== e.record.id)
    }
  })
})

onDestroy(() => {
  pb.realtime.unsubscribe('serverLogs')
})

function refresh() {
  pb.collection('serverLogs').getList(1, 200, {sort: '-created', filter: `server="${server.id}"`})
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
  pb.collection('serverLogs').delete(id)
}

function deletePastMessages() {
  // delete all messages except the last one
  logs.slice(1).forEach(l => pb.collection('serverLogs').delete(l.id))
}

function trustHostKey() {
  pb.collection('servers').update(server.id, {hostKey: details.payload})
  deleteMessage(details.id)
}

function trustHostName() {
  pb.collection('servers').update(server.id, {hostName: details.payload})
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

  <DataTable zebra title="Sync logs for {server.name}" description="Please wait for logs..." {headers} rows={logs}>
    <svelte:fragment slot="cell" let:row let:cell>
      {#if cell.key === "message"}
      <button type="button" class="bx--btn bx--btn--ghost" on:click={() => showDetails(row)}><Tag>{row.type}</Tag> {cell.value}</button>
      {:else if cell.key === "ago"}
        <time datetime={formatDateTime(row.created)}>{timeago.format(new Date(formatDateTime(row.created)))}</time>
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

<style>
  button.bx--btn--ghost {
    max-width: initial;
    padding: initial;
  }
</style>
