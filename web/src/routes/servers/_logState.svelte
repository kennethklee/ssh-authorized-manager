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

{#if row.state || LOG_TO_TAG_TYPE[row.state]}
  <Tag interactive type={LOG_TO_TAG_TYPE[row.state]} on:click={() => routeTo(`/servers/${row.id}/logs`)}>{LOG_TO_TAG_TEXT[row.state]}</Tag>
{:else}
  <Tag type="gray">Unknown</Tag>
{/if}
