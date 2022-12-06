<!--<script context="module" lang="ts">-->

<!--</script>-->
<script>
    import '../app.css'
    import { setContextClient } from '@urql/svelte'
    import { createGraphQLClient } from "$lib/http/index.js";
    import { base } from '$app/paths'
    console.log('using base', base)
    let client = createGraphQLClient()
    setContextClient(client)


</script>

<div class="min-h-screen flex flex-col">
    <div class='bg-blue-700 px-8 flex items-center justify-between py-4 shadow-sm text-white'> App Bar </div>
    <div class='relative flex flex-grow'>
        <nav class='bg-white shadow-sm p-6 space-y-6 w-64'> Navbar </nav>
        <main class='bg-gray-100 flex-1 p-6'><slot></slot></main>
    </div>
    <footer class="mt-auto">
        My Footer
    </footer>
</div>
