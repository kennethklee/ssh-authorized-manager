import PocketBase, { ClientResponseError } from 'pocketbase'

const client = new PocketBase()

client.afterSend = function (response, data) {
  if (!response.ok) {
    throw new ClientResponseError({
      url:      response.url,
      status:   response.status,
      data:     data,
    });
  }

  return data
}

/**
 * @param collectionName {string}
 * @param record {object}
 */
client.save = function (collectionName, record) {
  if (record.id) {
    return client.collection(collectionName).update(record.id, record)
  } else {
    return client.collection(collectionName).create(record)
  }
}

export default client
