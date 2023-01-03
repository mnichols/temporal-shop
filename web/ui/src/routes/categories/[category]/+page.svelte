<script lang="ts">
    import { PUBLIC_GRAPHQL_URL} from "$env/static/public";
    import { queryStore, gql, getContextClient } from '@urql/svelte';
    import { onMount} from "svelte";
    import {base} from "$app/paths";
    import Link from "$components/link/Link.svelte";
    import Default from "$components/layouts/Default.svelte";
    const inv = queryStore({
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
    import type { Inventory} from "../../gql/graphql";
    function toCategories(inventory: Inventory): Array<String> {
        return Array.from(new Set(inventory.games.map((g, i, games)=> g.category))).sort()
    }
    console.log('inv', inv)


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

        <h3>Categories</h3>
<!--        <div>{ JSON.stringify($inv.data.inventory.games) }</div>-->
        {#if $inv.fetching}
            <p>Loading...</p>
        {:else if $inv.error}
            <p>Error! {$inv.error.message}</p>
        {:else}
            {#each toCategories($inv.data.inventory) as category}
                <div><a href="{base}/categories/{category}">{ category }</a></div>
            {/each}
        {/if}
    </div>
</Default>


