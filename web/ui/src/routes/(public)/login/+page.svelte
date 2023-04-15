<script lang="ts">
    import '../../../app.css'
    import Button  from "$lib/holocene/button.svelte";
    import { doLogin } from '$lib/stores/auth-user'
    let password: string;
    let email: string;
    let message: string;
    const _doLogin = async () => {
        return doLogin({ email, password })
    }


</script>
<div class="min-h-screen bg-primary text-offWhite flex flex-col text-xl items-center justify-center">
    <header class="text-offWhite justify-center align-center flex flex-col">
        <h1 class="text-center text-4xl">Temporal Shop</h1>
    </header>
    <div class="flex flex-col text-inherit">
        <label class="text-inherit flex flex-col m-4 w-96" for="email">
            <span>Email</span>
            <input class="text-primary p-2" name="email" type="email" bind:value={ email } />
        </label>
        <label class="text-inherit flex flex-col m-4 w-96" for="password">
            <span>Password</span>
            <input class="text-primary p-2" name="password" type="password" bind:value={ password } />
        </label>
        <Button class="self-center border-spaceGray" variant="secondary" on:click={ _doLogin }>Login</Button>
        { #if message }
            <p class="error message">{ message }</p>
        { /if }
    </div>
</div>

<style lang="postcss">
    button {
        background: green;
    }
</style>