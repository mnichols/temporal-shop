<script lang="ts">
    import { login } from '$lib/auth'
    let password: string;
    let email: string;
    let message: string;
    const doLogin = async () => {
        return login({ email, password })
    }

</script>

<h1>Login</h1>

<div id="login">
    <label for="email">
        <input name="email" type="email" id="email" bind:value={ email } />
    </label>
    <label for="password">
        <input name="password" type="password" id="password" bind:value={ password } />
    </label>
    <button on:click={ doLogin }>Login</button>
    { #if message }
        <p class="error message">{ message }</p>
    { /if }
</div>

<style lang="postcss">
    button {
        background: green;
    }
</style>