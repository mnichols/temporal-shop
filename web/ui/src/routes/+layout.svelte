<script context="module" lang="ts">
    import '../app.css';
</script>
<script>
    import { setContextClient } from '@urql/svelte'
    import { createGraphQLClient } from "$lib/http/index.js";

    let client = createGraphQLClient()
    setContextClient(client)


</script>

<div class="foo">
<slot></slot>
</div>