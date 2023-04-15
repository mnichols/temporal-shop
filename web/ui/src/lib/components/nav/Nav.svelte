<script lang="ts">
    import Link from "$components/link/Link.svelte"
    import {base} from "$app/paths"

</script>
<h3><Link href="{base}">Home</Link></h3>
<ul>
    <li><Link href="{base}/shop/games">games</Link></li>
    <li><Link href="{base}/login">login</Link></li>
</ul>