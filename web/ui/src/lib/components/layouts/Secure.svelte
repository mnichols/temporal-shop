<script lang="ts">

    import { PUBLIC_GRAPHQL_URL} from "$env/static/public"
    import Nav from "$components/nav/Nav.svelte"
    import Default from "$components/layouts/Default.svelte"


</script>
<Default>
    <div slot="nav">
        <slot name="secure-nav"></slot>
        <Nav></Nav>
    </div>
    <div slot="main" class="text-offWhite">

        <slot name="secure-main"></slot>

        <h1>Temporal Shop</h1>

        <p>Connecting to graphql at {PUBLIC_GRAPHQL_URL}</p>

    </div>
</Default>
