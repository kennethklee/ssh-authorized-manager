import { writable } from 'svelte/store'
import pb from './pocketbase'

export const pathname = writable('')    // url path
export const params = writable({})      // url params i.e. `server` in /servers/:server
export const query = writable(new URLSearchParams())       // url query i.e. `q` in /search?q=foo

export const user = writable(null, function (set) {
  pb.authStore.model && getUserModel(pb.authStore.model).then(set)
  pb.authStore.onChange((_, model) => getUserModel(model).then(set))
  return () => { }
})

async function getUserModel(model) {
  if (!model || !model.id) return null; // logged out
  if ('collectionName' in model) {
    model.isUser = true
    return model; // logged in as normal user
  }

  // admin needs a user, so fetch
  console.log('admin needs a user, so fetch')
  return pb.collection('users').getFirstListItem(`email="${model.email}"`)
    .then(user => {
      console.log('got user', user)
      user.isAdmin = true
      user.isUser = true
      return user
    }).catch(err => {
      console.log('user not found', err)
      // user doesn't exist
      model.isAdmin = true
      model.isUser = false // cannot have servers
      return model
    })
}
