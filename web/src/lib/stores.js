import { writable } from 'svelte/store'
import pb from './pocketbase'

export const pathname = writable('')    // url path
export const params = writable({})      // url params i.e. `server` in /servers/:server
export const query = writable(new URLSearchParams())       // url query i.e. `q` in /search?q=foo

export const user = writable(null, function (set) {
  console.log('user store init')
  pb.authStore.model && getUserModel(pb.authStore.model).then(set)
  pb.authStore.onChange(async (_, model) => set(await getUserModel(model)))

  return () => { }
})

async function getUserModel(model) {
  if (!model || !model.id) return null; // logged out
  if (!('avatar' in model)) {
    model.isUser = true
    return model; // logged in as normal user (users don't have an avatar)
  }

  // admin needs a user, so fetch
  var users = await pb.users.getList(1, 1, { filter: `email="${model.email}"` })
  if (users.items.length) {
    // user exists
    users.items[0].isAdmin = true
    users.items[0].isUser = true  // has servers -- user account can store server relations
    return users.items[0]
  } else {
    // user doesn't exist
    model.isAdmin = true
    model.isUser = false // cannot have servers
    return model
  }
}
