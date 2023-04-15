<script lang="ts">
    import Default from "$components/layouts/Default.svelte"
    import Nav from "$components/nav/Nav.svelte";
</script>

<Default>
    <div slot="nav">
        <Nav />
    </div>
    <div slot="main" class="relative flex flex-col">
        <div class="flex flex-col flex-grow">
            <header class="relative m-0 p-0">
                <h2 class="flex flex-row mw-100 text-2xl justify-center flex-items-middle text-xl text-offWhite"><span>Shop</span></h2>
                <h3 class="absolute top-0 right-0 z-30 ">
                </h3>
            </header>
        </div>
    </div>
</Default>