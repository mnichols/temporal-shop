<script lang="ts">
    import { PUBLIC_GRAPHQL_URL} from "$env/static/public";
    import { queryStore, gql, getContextClient } from '@urql/svelte';
    import { onMount} from "svelte";
    import {base} from "$app/paths";
    import Link from "$components/link/Link.svelte";
    import Default from "$components/layouts/Default.svelte";
    const games = queryStore({
        client: getContextClient(),
        query: gql`
      query {
        inventory {
            games {
                id,
                product,
                image_url,
                category,
            }
        }
      }
    `,
    })

    console.log('games', games)


</script>
<Default>
    <div slot="nav">
        <h3>Routes</h3>
        <ul>
            <li><Link href="{base}/login">login</Link></li>
            <li><a href="{base}/foo">foo</a></li>
            <li><a href="{base}/foo/deep-foo">deep foo</a></li>
            <li><a href="{base}/login">login</a></li>
        </ul>
    </div>
    <div slot="main">
        <h1>Temporal Shop</h1>

        <p>Connecting to graphql at {PUBLIC_GRAPHQL_URL}</p>

    </div>
</Default>


