<script>
  import {
    Content,
    Grid,
    Row,
    Column,
    TooltipDefinition,
    FluidForm,
    TextInput,
    PasswordInput,
    Button,
    InlineNotification,
  } from 'carbon-components-svelte'
  import routeTo from 'page'
  import {onMount} from 'svelte'
  import pb from '$lib/pocketbase.js'


  var email
  var password
  var notify = null


  // check if already logged in
  onMount(() => {
    fetch('/api/me')
      .then(res => res.json())
      .then(data => {
        if (data) {
          routeTo('/')
        }
      })
  })

  function handleLogin(ev) {
    ev.preventDefault()
    notify = ''

    // check if admin
    var redirectPath = localStorage.getItem('redirectTo') ? localStorage.getItem('redirectTo') : '/'
    pb.admins.authViaEmail(email, password)
      .then(() => routeTo(redirectPath))
      .catch(() => {
        // check if user
        return pb.users.authViaEmail(email, password)
      })
      .then(() => routeTo(redirectPath))
      .catch(err => {
        notify = {kind: 'error', title: 'Error', subtitle: 'invalid email or password: ' + err.message}
      })
  }
</script>

<Content>
  {#if notify}
  <InlineNotification {...notify} />
  {/if}

  <Grid>
    <Row>
      <Column>
        <header>
          <TooltipDefinition style="float: right" tooltipText="Ask your admin to create an account for you. Admins: make sure to make yourself a user account.">Register</TooltipDefinition>
          <h1>Login</h1>
        </header>

        <FluidForm on:submit={handleLogin}>
          <TextInput required type="email" labelText="Email" placeholder="my-email@company.com..." autocomplete="email" bind:value={email} />
          <PasswordInput required labelText="Password" placeholder="********" autocomplete="current-password" bind:value={password} />

          <Button type="submit">Submit</Button>
        </FluidForm>
      </Column>
    </Row>
  </Grid>
</Content>

