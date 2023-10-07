<script>
import {
  Content,
  InlineNotification,
  Button,
  Form,
  FormGroup,
  TextInput,
  NumberInput,
  Toggle,
  PasswordInput,
  Loading,
  Modal,
  TextArea,
} from 'carbon-components-svelte'
import {onMount} from 'svelte'
import routeTo from 'page'
import pb from '$lib/pocketbase'


export var server = {port: 22, usePassword: true}
var title
var isConfirmDelete = false
var isSaving = false
var notify = null


onMount(() => {
  title = server.name
  document.title = title
})

function handleSave(ev) {
  ev.preventDefault()

  isSaving = true

  // Delete secret fields if not provided
  // if (server['password'] === '') {
  //   server = {...server, 'password': undefined}
  // }
  // if (server['privateKey'] === '') {
  //   server = {...server, 'privateKey': undefined}
  // }
  // if (server['privateKeyPassphrase'] === '') {
  //   server = {...server, 'privateKeyPassphrase': undefined}
  // }
  // replace all whitespace with blank
  // server.privateKey = server.privateKey.replaceAll(/\s/g, '')

  // @ts-ignore
  pb.save('servers', server)
    .then(data => {
      server = data
      title = server.name
      document.title = title
      notify = {kind: 'success', title: 'Saved', subtitle: 'Server saved. Check sync logs.'}
      routeTo(`/servers/${server.id}/logs`)
    })
    .catch(err => {
      if (err.data.data) {
        var fieldErrors = Object.keys(err.data.data).map(key => `${key}: ${err.data.data[key].message}`)
        notify = {kind: 'error', title: 'Error', subtitle: fieldErrors.join(', ')}
      } else {
        notify = {kind: 'error', title: 'Error', subtitle: err.message}
      }
    })
    .finally(() => isSaving = false)
}

function deleteServer(id) {
  isConfirmDelete = false
  pb.collection('servers').delete(id)
    .then(() => routeTo('/servers'))
}
</script>

<style>
  nav {
    display: flex;
    justify-content: space-between;
  }
</style>


{#if isSaving}
<Loading />
{/if}

<Modal
    bind:open={isConfirmDelete}
    danger
    modalHeading="Are you sure?"
    primaryButtonText="Delete"
    on:click:button--primary={() => deleteServer(server.id)}
    on:click:button--secondary={() => isConfirmDelete = false}>
    <p>This will delete the server and all associated logs.</p>
</Modal>

<Content>
  {#if notify}
  <InlineNotification {...notify} />
  {/if}

  <Form on:submit={handleSave}>
    <nav>
      <h2>{server.id ? title : 'New Server'}</h2>

      <article>
        {#if server.id}
        <Button kind="secondary" href="/servers/{server.id}/logs">Sync Logs</Button>
        <Button kind="danger" on:click={() => isConfirmDelete = true}>Delete</Button>
        {/if}
        <Button type="submit">Save</Button>
      </article>
    </nav>

    <article>
      <FormGroup><TextInput labelText="Name" placeholder="i.e. my host" required bind:value={server.name} /></FormGroup>
      <FormGroup><TextInput labelText="Host" placeholder="i.e. 127.0.0.1" required bind:value={server.host} /></FormGroup>
      <FormGroup><NumberInput label="Port" placeholder="i.e. 22" required bind:value={server.port} /></FormGroup>
      <FormGroup><TextInput labelText="Username" placeholder="i.e. myuser" required bind:value={server.username} /></FormGroup>

      <h4>Authentication</h4>
      <FormGroup><Toggle labelText="Use Password?" bind:toggled={server.usePassword}/></FormGroup>
      {#if server.usePassword}
      <FormGroup><PasswordInput labelText="Password" placeholder="********" bind:value={server['password']} /></FormGroup>
      {:else}
      <FormGroup><TextArea labelText="Private Key" bind:value={server['privateKey']} /></FormGroup>
      <FormGroup><PasswordInput labelText="Passphrase" bind:value={server['privateKeyPassPhrase']} /></FormGroup>
      {/if}

      <h4>Verification</h4>
      <p>Leave these blank if unknown. You will be asked to verify them in the sync logs if validation fails.</p>
      <FormGroup><TextInput labelText="Long host name" helperText="Found with `hostname -f` on the server" bind:value={server.hostname} /></FormGroup>
      <FormGroup><TextInput labelText="Host Key" bind:value={server.hostKey} /></FormGroup>

      <Button type="submit">Save</Button>
    </article>
  </Form>
</Content>
