<script lang="ts">
    import Default from "$components/layouts/Default.svelte"
    import Game from "$components/game/Game.svelte"
    import Button  from "$lib/holocene/button.svelte"
    import Nav from "$components/nav/Nav.svelte"
    import { page } from '$app/stores'
    import { goto } from '$app/navigation'
    import {getContextClient, queryStore} from "@urql/svelte";
    import CartButton from "$lib/components/cart/CartButton.svelte"
    import {InventoryDocument} from '$gql'
    import type {  InventoryInput} from '$gql'

    /** @type {import('../../../../../.svelte-kit/types/src/routes').PageData} */
    export let data;
    let input: InventoryInput = { category: data.category }
    page.subscribe(ignore => {
        input.category = $page.url.searchParams.get('category')
    })

    $: console.log('page category be ' + input)
    $: inventory = queryStore({
        client: getContextClient(),
        query: InventoryDocument,
        variables: { input },
    })
    // event handler to filter listing by category
    function filterByCategory(category: String) {
        const newUrl = new URL($page.url);
        newUrl?.searchParams?.set('category', category);
        console.log('filtering', category, 'url', newUrl)

        return goto(newUrl, { invalidateAll: true});
    }
</script>
<Default>
    <div slot="nav">
        <Nav />
    </div>
    <div slot="main" class="relative flex flex-col">
        <div class="flex flex-col flex-grow">
            <header class="relative m-0 p-0">
                <h2 class="flex flex-row mw-100 text-2xl justify-center flex-items-middle text-xl text-offWhite"><span>Games</span></h2>
                <h3 class="absolute top-0 right-0 z-30 ">
                    <CartButton></CartButton>
                </h3>
            </header>

            {#if $inventory}
                {#if $inventory.loading}
                    <p>Loading...</p>
                {:else if $inventory.errors}
                    {#each $inventory?.errors as error}
                        <p>Error! {$inventory.error.message}</p>
                        {/each}
                {:else}
                    {#if $inventory?.data?.inventory?.categories}
                        <div class="flex flex-row flex-start sticky top-0 z-30 bg-offWhite m-2 border-0">
                        {#each $inventory?.data?.inventory?.categories as cat}
                            <Button class="m-2 flex-initial justify-center flex" on:click={(e)=> filterByCategory(cat)} variant={cat==input.category ? 'primary':'secondary'}>{cat}</Button>
                        {/each}
                        </div>
                        {/if}
                    {#if $inventory?.data?.inventory?.games}
                    <div class="relative flex flex-row flex-wrap">
                        {#each $inventory?.data?.inventory?.games as game}
                            <Game game={game}></Game>
                        {/each}
                    </div>
                        {/if}
                {/if}
            {/if}
        </div>
    </div>
</Default>


