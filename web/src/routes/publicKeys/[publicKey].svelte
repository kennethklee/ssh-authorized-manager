<script>
import {
  Content,
  InlineNotification,
  Button,
  Form,
  FormGroup,
  TextInput,
  TextArea,
  Loading,
  Modal,
} from 'carbon-components-svelte'
import routeTo from 'page'
import pb from '$lib/pocketbase.js'
import {user} from '$lib/stores'

export var publicKey = {}
var notify = null
var title
var isSaving = false
var isConfirmDelete = false


$: initOnce(publicKey)  // initialize server with item once

function initOnce(publicKey) {
  title = publicKey.comment || (publicKey.publicKey ? publicKey.publicKey.slice(0, 10) + '...' : 'New Public Key')
  document.title = title
}

function handleSave(ev) {
  ev.preventDefault()

  isSaving = true
  publicKey.userId = publicKey.userId || $user.id

  // @ts-ignore
  pb.save('publicKeys', publicKey)
    .then(data => {
      publicKey = data
      isSaving = false
      notify = {kind: 'success', title: 'Saved', message: 'Public key saved'}

      routeTo(`/publicKeys/${publicKey.id}`)
    })
    // TODO handle validation errors
    .catch(err => notify = {kind: 'error', title: 'Error', subtitle: err.message})
    .finally(() => isSaving = false)
}

function handlePublicKeyChange(ev) {
  // format <type> <key> <comment>
  var parts = publicKey.publicKey.trim().split(' ')
  if (parts.length === 3) {
    // Usually <type> <key> <comment>
    publicKey.type = parts[0]
    publicKey.publicKey = parts[1]
    publicKey.comment = parts[2]
  } else if (parts.length === 2) {
    // Usually <type> <key>
    publicKey.type = parts[0]
    publicKey.publicKey = parts[1]
  }
}

function deletePublicKey(id) {
  pb.collection('publicKeys').delete(id)
    .then(() => routeTo('/publicKeys'))
    .catch(err => notify = {kind: 'error', title: 'Error', subtitle: err.message})
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
    on:click:button--primary={() => deletePublicKey(publicKey.id)}
    on:click:button--secondary={() => isConfirmDelete = false}>
    <p>This will delete the public key.</p>
</Modal>

<Content>
  {#if notify}
  <InlineNotification {...notify} />
  {/if}

  <nav>
    <h2>{title}</h2>
    <article>
      {#if publicKey.id}
      <Button kind="danger" on:click={() => isConfirmDelete = true}>Delete</Button>
      {/if}
      <Button type="submit">Save</Button>
    </article>
  </nav>

  <Form title="New" on:submit={handleSave}>
    <FormGroup><TextArea labelText="Public Key" required placeholder="my public key..." bind:value={publicKey.publicKey} on:change={handlePublicKeyChange}/></FormGroup>
    <FormGroup><TextInput labelText="Type" placeholder="ed25519..." bind:value={publicKey.type} /></FormGroup>
    <FormGroup><TextInput labelText="Name" placeholder="my computer..." bind:value={publicKey.comment} /></FormGroup>

    <Button type="submit">Save</Button>
  </Form>
</Content>
