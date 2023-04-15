<script lang="ts">
    import Default from "$components/layouts/Default.svelte"
    import Button  from "$lib/holocene/button.svelte"
    import Nav from "$components/nav/Nav.svelte"
    import {getContextCart} from "$lib/stores/cart/cart"

    import SelectableTable from '$lib/holocene/table/selectable-table.svelte'
    import SelectableTableRow from '$lib/holocene/table/selectable-table-row.svelte'
    import BulkActionButton from '$lib/holocene/table/bulk-action-button.svelte'

    import Table from '$lib/holocene/table/table.svelte'
    import TableRow from '$lib/holocene/table/table-row.svelte'
    import Select from '$lib/holocene/select/select.svelte'
    import Option from '$lib/holocene/select/option.svelte'

    import { loadStripe } from '@stripe/stripe-js'
    import { Elements } from 'svelte-stripe'
    import {onDestroy, onMount} from 'svelte'

    let stripe = null
    let cartStore = getContextCart()
    $: cart = cartStore
    onMount(async () => {
        // stripe = await loadStripe(PUBLIC_STRIPE_KEY)
    })

    let selectedItems = []
    interface Opt {
        value: number
        selected: boolean
    }
    const MinCount = 10
    const listQuantity = (current: number,count: number = 0):[Opt]  =>  {
        count = Math.max(count, MinCount)

        return new Array(count)
            .fill(0)
            .map((a, i) =>{ return { selected: (i + 1) === current, value: i+1}})
    }

    const setQuantity = ({productId}, newQuantity) => {
        cartStore.setItemQuantity({productId, quantity: newQuantity})
    }
    const removeItems = (e) => {
        let args = selectedItems.map(i => {
            return { productId: i.productId, quantity: 0 }
        })
        cartStore.putItems(args)
    }
</script>
<style>
    .currency::before {
        content:"$";
    }
    .percent::before {
        content:"%";
    }
</style>
<Default>
    <div slot="nav">
        <Nav />
    </div>
    <div slot="main" class="relative flex flex-col">
        <div class="flex flex-col flex-grow">
            <header class="relative m-0 p-0">
                <h2 class="flex flex-row mw-100 text-2xl justify-center flex-items-middle text-xl text-offWhite"><span>Cart</span></h2>
            </header>

            {#if $cart}
                {#if $cart.fetching}
                    <p>Loading...</p>
                {:else if $cart.errors}
                    {#each $cart?.errors as error}
                        <p>Error! {$cart.error.message}</p>
                    {/each}
                {:else}
                    <div class="flex flex-row">
                        {#if $cart?.data?.cart?.items}
                            <SelectableTable class="text-offWhite w-3/6" items={$cart.data.cart.items} bind:selectedItems>
                                <svelte:fragment slot="bulk-action-headers">
                                    <th>
                                        {selectedItems.length} Selected
                                    </th>
                                    <th class="w-1/6">
                                        <BulkActionButton variant="destructive" on:click={removeItems}>Remove</BulkActionButton>
                                    </th>
                                    <th class="w-64" />
                                </svelte:fragment>
                                <svelte:fragment slot="default-headers">
                                    <th>Title</th>
                                    <th class="w-1/6">Price</th>
                                    <th class="w-32">Quantity</th>
                                </svelte:fragment>
                                {#each $cart.data.cart.items as item, i}
                                    {#if item}
                                        <SelectableTableRow
                                                selected={selectedItems.includes(item)}
                                                item={item}
                                                class="text-primaryText bg-offWhite"
                                        >
                                            <td>{item.title}</td>
                                            <td>{item.price}</td>
                                            <td class="flex flex-row">
                                                <Select
                                                        id="qty-{item-i}"
                                                        label="quantity"
                                                        placeholder="quantity"
                                                        value={item.quantity}
                                                        onChange={(e) => setQuantity(item, e)}
                                                        class="w-14"
                                                >
                                                    {#each listQuantity(item.quantity) as qty}
                                                        <Option value={qty.value}>{qty.value}</Option>
                                                    {/each}

                                                </Select>
                                            </td>

                                        </SelectableTableRow>
                                        {/if}
                                {/each}
                            </SelectableTable>
                            <div class="bg-offWhite m-4 rounded-lg flex flex-col w-3/12">
                                <div class="text-primaryText p-4 pb-0 flex flex-row font-medium justify-between">
                                    <span>Subtotal:</span>
                                    <span class="currency">{$cart.data.cart.subtotal}</span>
                                </div>
                                <div class="text-primaryText p-4 pb-0 flex flex-row font-medium justify-between">
                                    <div class="flex flex-row">Tax<span class="percent">{$cart.data.cart.taxRate}</span>:</div>
                                    <span class="currency">{$cart.data.cart.tax}</span>
                                </div>
                                <div class="text-primaryText p-4 pb-0 flex flex-row font-bold justify-between">
                                    <span>Total Due:</span>
                                    <span class="currency">{$cart.data.cart.total}</span>
                                </div>
                                <div class="text-primaryText border-b-2 p-4 flex flex-row font-medium justify-between">
                                    <h2 class="text-primaryText border-b-2 p-4">Payment</h2>
                                    {#if stripe}
                                        <Elements {stripe}>
                                            <!-- this is where your Stripe components go -->
                                        </Elements>
                                    {/if}
                                </div>
                            </div>
                        {/if}
                    </div>
                {/if}
            {/if}
        </div>
    </div>
</Default>


