<script lang="ts">
    export let productId: String
    export let quantity = 1
    // import cart from "$lib/stores/cart/cart";
    import Add from "$lib/holocene/icon/svg/add.svelte";
    import { getContextCart} from "$lib/stores/cart/cart";
    let cart = getContextCart()
    async function onClick(e) {
        console.log('AddCartItem.svelte#onClick', e, 'productId',productId, 'quantity', quantity)
        await cart.addItemToCart({productId, quantity})

        console.log('AddCartItem.svelte#onClick DONE')
    }
</script>
<button class="absolute
        top-0
        right-0
        -mr-3
        -mt-3
        rounded-full
        border-1
        border-blue-500
        bg-blue-500
        text-offWhite
        text-2xl
        flex
        items-center
        justify-center
        w-10
        h-10
        p-0
" type="button" on:click={ onClick }>
    <span><Add></Add></span>
</button>