import { subscriptionStore} from '@urql/svelte'
import type {Client} from "@urql/core";
import {CartDocument} from "../../gql";
import {writable} from "svelte/store";

export const createThing = (client: Client) => {
    let ss
    let mystore = writable()
    mystore.subscribe(arg => {
        ss = subscriptionStore( {
            client,
            query: CartDocument,
            variables: {input:{cartId:''}},
        })

        ss.subscribe(arg2 => console.log('subscription store','original value', arg,'subvalue', arg2))
    })
    return mystore
}