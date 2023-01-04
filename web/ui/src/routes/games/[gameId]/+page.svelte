<script lang="ts">
    import { queryStore, gql, getContextClient } from '@urql/svelte';
    import {base} from "$app/paths";
    import Link from "$components/link/Link.svelte";
    import Default from "$components/layouts/Default.svelte";
    import type { InventoryInput} from "../../../gql/graphql";

    /** @type {import('./$types').PageData} */
    export let data;
    const input: InventoryInput = { categoryId: data.categoryId }
    const categoriesQuery = gql`
          query ($input: InventoryInput!) {
            inventory(input: $input) {
                games {
                    id,
                    product,
                    imageUrl,
                    category,
                }
            }
          }
        `
    $: inv = queryStore({
        client: getContextClient(),
        query: categoriesQuery,
        variables: {input},
    })
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

        <h2>{ data.categoryId}</h2>
        <h3>Games</h3>
            {#if $inv.fetching}
                <p>Loading...</p>
            {:else if $inv.error}
                <p>Error! {$inv.error.message}</p>
            {:else}
                {#each $inv.data.inventory.games as game}
                    <div>{ game.product }</div>
                {/each}
            {/if}
    </div>
</Default>


