<script>
  import {
    Content,
    SkeletonPlaceholder,
    Tile,
    CopyButton,
  } from "carbon-components-svelte";
  import pb from '$lib/pocketbase.js'
  import {sshTarget} from '$lib/utils.js'
</script>

<style>
  h2 {
    user-select: none;
  }
  .grid {
    display: grid;
    grid-gap: 1em;
    grid-template-columns: repeat( auto-fit, minmax(100px, 1fr) );
  }
  .host {
    display: flex;
    flex-direction: row;
    line-height: 40px;
  }
</style>

<!-- TODO quick links to recent items -->

<Content>
  <header>
    <h1>Welcome</h1>
  </header>

  <article class="grid">
    {#await pb.collection('servers').getList(1, 8, {sort: '-updated', filter: 'lastState="info"'})}
      <SkeletonPlaceholder style="width: auto" />
      <SkeletonPlaceholder style="width: auto" />
      <SkeletonPlaceholder style="width: auto" />
    {:then response}
      {#each response.items as server}
        <Tile>
          <h2>{server.name}</h2>
          <div class="host">
            <span>{sshTarget(server)}</span>
            <CopyButton text={sshTarget(server)} />
          </div>
        </Tile>
      {:else}
        <p>No servers found. Please ask an admin to add a server to your account.</p>
      {/each}
    {/await}
  </article>
</Content>
