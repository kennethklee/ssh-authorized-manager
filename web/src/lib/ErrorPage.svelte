<script>
  import {
    Content,
  } from 'carbon-components-svelte'

  export var code = 404
  export var message = 'Not Found'
</script>


<Content>
  <h2>{code}</h2>
  <h1>{message}</h1>
</Content>
