<script lang="ts">
    import Cart from "$lib/holocene/icon/svg/cart.svelte"
    import {getContextCart} from "$lib/stores/cart/cart"
    import type { CartStore } from '$lib/stores/cart/cart'
    import Badge from '$lib/holocene/badge.svelte'
    import {go} from "$lib/nav/index.js";
    let cart = getContextCart()
    $: itemCount = cart.itemCount

</script>

<a class="cursor-pointer flex flex-row relative justify-center items-center rounded-lg bg-offWhite border-offWhite border-1 text-primary h-8 w-16"
    on:click|preventDefault={ () => go('/shop/cart')}
>
    <Cart viewBox="0 0 1024 1024" ></Cart>
    <Badge class="badge absolute top-0 right-0 origin-bottom-left translate-y-[-10px] translate-x-[10px]"
            type="count">{$itemCount}</Badge
    >
</a>