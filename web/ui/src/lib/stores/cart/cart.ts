import type { Readable, Writable,} from 'svelte/store'
import {writable, get, derived,} from "svelte/store";
import type { SetCartItemsInput, CartItemInput} from '$gql';
import {queryStore, subscriptionStore, mutationStore,} from "@urql/svelte";
import type { OperationResultStore, Pausable ,OperationResultState,} from '@urql/svelte'
import {CartDocument, CurrentCartDocument, SetCartItemsDocument} from "$gql";
import { getContext, setContext } from 'svelte'
import type { Client } from '@urql/svelte'
import type {
    CartSubscription,
    CartSubscriptionVariables,
    CurrentCartQuery,
    CurrentCartQueryVariables
} from "$gql";
import { sortBy } from 'lodash'
import { Logger} from '$log'

const _contextKey = '$$_cart';

export type ComposedData = CurrentCartQuery | CartSubscription
export type ComposedVariables = CurrentCartQueryVariables | CartSubscriptionVariables

export type ComposedStore = OperationResultState<ComposedData, ComposedVariables>

type CartSubscriptionStore = OperationResultStore<CartSubscription, CartSubscriptionVariables>

export interface CartStore extends Readable<ComposedStore> {
    addItemToCart(args: CartItemInput): Promise<null>
    setItemQuantity(args: CartItemInput): Promise<null>;
    putItems(args: [CartItemInput]): Promise<null>
    itemCount: Readable<number>
}

type dedupedItems = {
    [productId: string]: CartItemInput;
}
export const getContextCart = (): CartStore => {
    const out = getContext(_contextKey);
    if (process.env.NODE_ENV !== 'production' && !out) {
        throw new Error(
            'No Cart was found in Svelte context. Did you forget to call setContextCart?'
        );
    }

    return out as CartStore;
}
export const setContextCart = (cart: CartStore): void => {
    setContext(_contextKey, cart)
}

const dedupeItems = (items:CartItemInput[]): dedupedItems => {
    let deduped: dedupedItems = {}
    return items.reduce((acc, item) => {
        let ref = acc[item.productId]
        if (ref) {
            ref.quantity += item.quantity
        } else {
            acc[item.productId] = item
        }
        return acc
    }, deduped)
}
// const setItemQuantity = writable<CartItemInput>()
// putItem unifies the PUT mutation of items, accepting a `items` filter before applying
const putItem = (client: Client, cart: ComposedData, args:CartItemInput,
                       filterItemsFn: (current: CartItemInput[]) => CartItemInput[]): null => {

    const { id, items } = cart.cart
    const input:SetCartItemsInput =
        {
            cartId: id,
            items: Object.values(dedupeItems(filterItemsFn(items.map(({productId, quantity}) => ({productId, quantity}))))),
        }
    try {
        Logger.debug(input,'#setCartItems')
        let m = mutationStore({
            client,
            query: SetCartItemsDocument,
            variables: { input},
        })
    } catch(e) {
        Logger.error(e, 'setCartItemsStore errored')
    }
    return null
}

const calcItemCount = (cart?: ComposedData): number => {
    if(!cart || !cart?.cart?.items) {
        return 0
    }
    let { items } = cart?.cart
    return (items || []).reduce((acc: number, {quantity} ) => {
        acc = acc + quantity
        return acc
    }, 0)
}
const withOrderedItems = (arg: ComposedStore): ComposedStore => {
    if(!arg?.data?.cart) {
        return arg
    }
    arg.data.cart.items = sortBy(arg.data.cart.items, ['category', 'title'])
    return arg
}

export const createCartStore = (client:Client): CartStore => {
    let out = writable<ComposedStore>()
    let qInput: CurrentCartQueryVariables = { input: undefined }

    let ss: CartSubscriptionStore
    let qs = queryStore({
        client,
        query: CurrentCartDocument,
        variables: qInput,
        requestPolicy: 'cache-and-network'
    })
    qs.subscribe(arg => {
        Logger.debug(arg, 'queryStore subscribe')
        if(arg.stale) {
            return
        }
        out.set(withOrderedItems(arg))

        if(!ss && arg?.data?.cart?.id) {
            ss = subscriptionStore({ client, query: CartDocument, variables: { input: { cartId: arg?.data?.cart?.id}}})
            ss.subscribe(arg2 => {

                Logger.debug(arg2, 'subscriptionStore %s','received update from cart subscription')
                if(arg2.stale || (!arg2.data && !arg2.error)) {
                    return
                }
                out.set(withOrderedItems(arg2))
                return () => {
                    Logger.debug('ss.unsubscribing')
                }
            })
        }
        return () => {
            Logger.debug('qs.unsubscribing')
        }
    })

    let itemCountStore = derived([out], ([$out]) => {
        if($out?.data?.cart?.items) {
            let result = calcItemCount($out?.data)
            return result
        }
        return 0
    })

    return {
        subscribe: out.subscribe,
        itemCount: itemCountStore,
        addItemToCart: async (args: CartItemInput): Promise<null> => {
            Logger.debug(args, '#addItemToCart')
            let cart = get(out)
            if(!cart?.data) {
                throw new Error('addItemToCart is not valid without a cart')
            }
            return putItem(client, cart.data, args, (current: CartItemInput[]) => {
                return [...current, args]
            })
        },
        setItemQuantity: async (args: CartItemInput): Promise<null> => {
            Logger.debug(args, '#setItemQuantity')
            let cart = get(out)
            if(!cart?.data) {
                throw new Error('addItemToCart is not valid without a cart')
            }
            return putItem(client, cart.data, args, (current: CartItemInput[]) => {
                return [...current.filter((cur: CartItemInput) => cur.productId !== args.productId), args]
            })
        },
        putItems: async (args: [CartItemInput]): Promise<null> => {
            Logger.debug(args, '#putItems')

            let cart = get(out)
            if(!cart?.data) {
                throw new Error('putItems is not valid without a cart')
            }
            let pids = args.reduce((cur, item) => {
                cur[item.productId] = item
                return cur
            }, {})
            return putItem(client, cart.data, args, (current: CartItemInput[]) => {
                return [...current.filter((cur: CartItemInput) => !!!pids[cur.productId]), ...args]
            })
        },
    }
}