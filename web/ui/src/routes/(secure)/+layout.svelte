<!--<script context="module" lang="ts">-->

<!--</script>-->
<script lang="ts">
    import {setContextClient} from "@urql/svelte";
    import { createClient } from '$lib/http/urql'
    import '../../app.css'
    import {createUserStore, goLogin, setContextUser} from "$lib/stores/auth-user";
    const client = createClient()
    setContextClient(client)
    let user = createUserStore(client)
    setContextUser(user)
    $: if(!user) {
        goLogin()
    }

</script>
<slot></slot>
<!--<div class="min-h-screen flex flex-col">-->
<!--    <header class='sticky top-0 z-30 bg-primary px-8 flex items-center justify-between py-4 shadow-sm text-offWhite'>-->
<!--        <h1 class="text-lg">Temporal Shop</h1>-->
<!--    </header>-->
<!--    <div class='relative flex flex-grow'>-->
<!--        <nav class='bg-secondary shadow-sm p-6 space-y-6 w-64 text-offWhite'> <slot name="nav"></slot></nav>-->
<!--        <main class='bg-spaceGray flex-1 p-6'><slot name="main"></slot></main>-->
<!--    </div>-->
<!--    <footer class="mt-auto">-->
<!--        My Footer-->
<!--    </footer>-->
<!--</div>-->