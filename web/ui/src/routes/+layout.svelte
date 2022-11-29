<script context="module" lang="ts">
    import '../app.css';
</script>
<script>
    import { createClient, setContextClient } from '@urql/svelte'
    import { PUBLIC_GRAPHQL_URL} from "$env/static/public";
    import {withSecurityOptions} from "$lib/http/index.js";

    const client = createClient({
        url: PUBLIC_GRAPHQL_URL,
        fetchOptions: withSecurityOptions({}, true),
    });

    setContextClient(client);
</script>

<div class="foo">
<slot></slot>
</div>