<script lang="ts">
    import { PUBLIC_GRAPHQL_URL} from "$env/static/public";
    import { queryStore, gql, getContextClient } from '@urql/svelte';
    import { onMount} from "svelte";
    import {base} from "$app/paths";
    import Link from "$components/link/Link.svelte";

        const games = queryStore({
            client: getContextClient(),
            query: gql`
      query {
        games {
            items {
                id,
                product,
                image_url,
                category,

            }
        }
      }
    `,
        })

</script>
<h1>Temporal Shop</h1>


<p>Connecting to graphql at {PUBLIC_GRAPHQL_URL}</p>
<h3>Routes</h3>
<ul>
    <li><Link href="{base}/login">login</Link></li>
    <li><a href="{base}/foo">foo</a></li>
    <li><a href="{base}/foo/deep-foo">deep foo</a></li>
    <li><a href="{base}/login">login</a></li>
</ul>