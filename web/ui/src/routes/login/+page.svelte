<script lang="ts">
    import { goto } from "$app/navigation";
    let password: string;
    let email: string;
    let message: string;
    async function login() {
        try {
            const user: UserCredentials = await auth.signIn({ email, password });
            if (user) {
                goto("/");
            }
        } catch(error) {
            message = error;
        }
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
    <button on:click={ login }>Login</button>
    { #if message }
        <p class="error message">{ message }</p>
    { /if }
</div>

<style lang="postcss">
    p {
        background: turquoise;
    }
    p {
        color: white;
    }
</style>