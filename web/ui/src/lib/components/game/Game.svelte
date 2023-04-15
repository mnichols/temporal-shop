<script lang="ts">
    import AddCartItem from "$lib/components/add-cart-item/AddCartItem.svelte";
    export let game
</script>
<article class="relative bg-offWhite relative border-2 border-gray-400 basis-3/12 flex flex-col m-4 p-2 rounded-md">

    <h3 class="font-bold text-lg text-center justify-center h-18 grow-2 p-2">{ game.id }</h3>
    <figure class="relative flex flex-col basis-full justify-center">
        <img class="basis-full" src="{ game.imageUrl }" alt="game pic" />
        <figcaption class="mt-auto text-center font-semibold">${ game.price}</figcaption>
    </figure>
    <AddCartItem productId={ game.id}></AddCartItem>
</article>