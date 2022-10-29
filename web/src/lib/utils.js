/**
 * Formats date time to cross-compatible format
 * @param {string} datetime Date time string YYYY-MM-DD HH:MM:SS.sss
 * @returns {string} Formatted date time string YYYY-MM-DDTHH:MM:SS.sssZ
 */
export function formatDateTime(datetime) {
  // converts YYYY-MM-DD HH:MM:SS.sss to YYYY-MM-DDTHH:MM:SS.sssZ
  return datetime.trim().replace(' ', 'T') + 'Z'
}

export function sshTarget(server) {
  return `${server.username}@${server.host}${server.port = 22 ? '' : ':' + server.port}`
}
