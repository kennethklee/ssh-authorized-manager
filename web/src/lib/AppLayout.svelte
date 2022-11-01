<script>
  import {
    Header,
    HeaderUtilities,
    HeaderAction,
    HeaderPanelLinks,
    HeaderPanelDivider,
    HeaderPanelLink,
    SideNav,
    SideNavMenu,
    SideNavItems,
    SideNavLink,
    SideNavDivider,
    SkipToContent,
  } from "carbon-components-svelte";
  import routeTo from 'page'
  import pb from "$lib/pocketbase.js"
  import {user} from '$lib/stores.js'

  import DashboardReference from "carbon-icons-svelte/lib/DashboardReference.svelte";
  import UserAvatarFilledAlt from "carbon-icons-svelte/lib/UserAvatarFilledAlt.svelte";

  let isSideNavOpen = false;
  let isOpen1 = false;

  function handleLogOut(ev) {
    ev.preventDefault();
    pb.authStore.clear();
    routeTo('/login');
  }
</script>

<style>
  li {
    padding: .375rem 1rem;
  }
</style>

<Header platformName="SSH Authorized Manager" href="/" bind:isSideNavOpen>
  <svelte:fragment slot="skip-to-content">
    <SkipToContent />
  </svelte:fragment>

  <HeaderUtilities>
    {#if $user}
    <span>{$user.name || $user.email}</span>
    <HeaderAction bind:isOpen={isOpen1} icon={UserAvatarFilledAlt} closeIcon={UserAvatarFilledAlt}>
      <HeaderPanelLinks>
        <HeaderPanelDivider>Account</HeaderPanelDivider>
        <li>{$user?.profile ? $user.profile.name : $user?.email}</li>
        <HeaderPanelLink on:click={handleLogOut}>Log Out</HeaderPanelLink>

        {#if $user?.isAdmin}
        <HeaderPanelDivider>Admin</HeaderPanelDivider>
        <HeaderPanelLink href="/_" target="_blank">Dashboard</HeaderPanelLink>
        {/if}

        <HeaderPanelDivider>Version</HeaderPanelDivider>
        <HeaderPanelLink>{import.meta.env.APP_VERSION}</HeaderPanelLink>
      </HeaderPanelLinks>
    </HeaderAction>
    {/if}
  </HeaderUtilities>
</Header>

{#if $user}
<SideNav bind:isOpen={isSideNavOpen}>
  <SideNavItems>
    {#if $user?.isUser}
    <SideNavLink text="My Public Keys" href="/publicKeys" />
    <SideNavLink text="My Servers" href="/servers" />
    {:else}
    <p>To save your public keys and servers, <a href="/_/#/users">create a user for yourself</a>.</p>
    {/if}

    <!-- Admin -->
    {#if $user?.isAdmin}
    <SideNavDivider />
    <SideNavMenu text="Manage" open>
      <SideNavLink text="Servers" href="/servers?all" />
      <SideNavLink text="Users" href="/users" />
    </SideNavMenu>
    {/if}
  </SideNavItems>
</SideNav>
{/if}

<slot />
