<script>
import {onMount, tick} from 'svelte'
import router from 'page'

import {pathname, params, query, user} from '$lib/stores.js'
import pb from '$lib/pocketbase.js'
import AppLayout from '$lib/AppLayout.svelte'
import ErrorPage from '$app/lib/ErrorPage.svelte'


var pageComponent = null
var pageArgs = {}

function goto404(err) {
  console.error(err)
  pageArgs = {message: 'Not Found', code: 404}
  pageComponent = ErrorPage
}

// Middleware
router('*', async (ctx, next) => {
  // console.log(ctx)
  pageComponent = null  // Reset page component so we don't leak state
  await tick()
  pageArgs = {}
  query.update(q => new URLSearchParams(ctx.querystring))
  params.update(p => ctx.params)
  pathname.update(p => ctx.pathname)
  next()
})

// Build routes statically
var pageComponentMap = import.meta.glob('./routes/**/*.svelte', {eager: true})
Object.keys(pageComponentMap).forEach(pageComponentKey => {
  var chainCallbacks = []
  var routeParts = []
  var parts = pageComponentKey.split('/')
  for (var i = 2; i < parts.length; i++) {
    var part = parts[i]
    // skip component if starts with `_`
    if (part.startsWith('_')) return;

    // strip out the .svelte
    if (part.endsWith('.svelte')) {
      part = part.slice(0, -7)
    }
    // converts [param] to :param
    if (part.startsWith('[') && part.endsWith(']')) {
      part = ':' + part.slice(1, -1)
    }

    // attempt to load the param
    if (part.startsWith(':')) {
      var param = part.slice(1)
      if (part === ':user') {
        chainCallbacks.push(async (ctx, next) => {
          console.debug('load user', ctx.params.user)
          pb.collection('users').getOne(ctx.params.user)
            .then(data => pageArgs.user = data)
            .then(next)
            .catch(goto404)
        })
      } else {
        // luckily, our collection names are it's singular form + 's', and params are it's singular form.
        var collectionName = param + 's'  // singular form (param) + 's'
        chainCallbacks.push(async (ctx, next) => {
          if (ctx.params[param] === 'new') return next();
          console.debug('load', param, ctx.params[param])
          pb.collection(collectionName).getOne(ctx.params[param])
            .then(data => pageArgs[param] = data)
            .then(next)
            .catch(goto404)
        })
      }
    }

    if (part !== 'index') {
      routeParts.push(part)
    }
  }

  // @ts-ignore
  chainCallbacks.push(ctx => pageComponent = pageComponentMap[pageComponentKey].default)
  router('/' + routeParts.join('/'), ...chainCallbacks)
  // console.log('/' + routeParts.join('/'))   // uncomment to show all routes
})

// Any other route
router('/*', ctx => {
  pageArgs = {message: 'Not Found', code: 404}
  pageComponent = ErrorPage
})

// Start router
onMount(() => {
  // Check for auth. Also avoid redirect loop if we're already on the login page.
  if (!pb.authStore.isValid && !$user && !window.location.pathname.startsWith('/login')) {
    // check /api/me for header auth (already logged in)
    fetch('/api/me')
      .then(res => res.json())
      .then(data => {
        if (data) {
          // Login to refresh authStore with JWT token (pocketbase SDK needs this)
          if (data.profile) {
            pb.collection('user').authWithPassword('dontmatter', 'dontcare')
          } else {
            pb.admins.authWithPassword('dontmatter', 'dontcare')
          }
          router.start()
        } else {
          // save `redirectTo`
          localStorage.setItem('redirectTo', window.location.pathname)

          router.start()
          router('/login')
        }
      })

  } else {
    router.start()
  }
})
</script>

<AppLayout>
  <svelte:component this={pageComponent} {...pageArgs} />
</AppLayout>
