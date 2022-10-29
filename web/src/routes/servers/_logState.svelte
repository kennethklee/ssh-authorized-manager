<script>
import {
  Tag,
} from 'carbon-components-svelte'
import routeTo from 'page'

const LOG_TO_TAG_TYPE = {
  'error': 'red',
  'info': 'green',
  'hostKey': 'purple',
  'hostName': 'purple',
}
const LOG_TO_TAG_TEXT = {
  'error': 'ERROR',
  'info': 'OK',
  'hostKey': 'ACTION NEEDED',
  'hostName': 'ACTION NEEDED',
}

export var row = {}
</script>

{#if row.lastLog || LOG_TO_TAG_TYPE[row.lastLog?.type]}
  <Tag interactive type={LOG_TO_TAG_TYPE[row.lastLog?.type]} on:click={() => routeTo(`/servers/${row.id}/logs`)}>{LOG_TO_TAG_TEXT[row.lastLog?.type]}</Tag>
{:else}
  <Tag type="gray">Unknown</Tag>
{/if}
