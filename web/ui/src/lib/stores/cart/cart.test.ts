import {CartStore, ComposedStore, createCartStore} from './cart'
import {afterEach, beforeAll, beforeEach} from "vitest"
import {get, writable} from "svelte/store";
import {temporalShop, use} from "$lib/http/mock-server";
import { find, isEqual } from 'lodash'
import {fail} from "@sveltejs/kit";
import { rest } from 'msw'
import chaiExclude from 'chai-exclude'
import type {Client} from "@urql/core";
import {createTestClient} from "../../http/urql";
import {CartDocument} from "../../../gql";
import type { SetCartItemsInput} from '$gql'
chai.use(chaiExclude)
import { subscriptionStore} from "@urql/svelte";

describe('createStore', async () => {

    const initialCart = {
        cart: {
            id: 'cart-1',
            shopperId: 'shopper-1',
            subtotal: '19.99',
            total: '23.99',
            tax: '4.00',
            taxRate: '20',
            timestamp: '1970-04-05',
            items: [
                {
                    productId: 'p1',
                    quantity: 1,
                    subtotal: '19.99',
                    price: '19.99',
                    title: 'p1'
                }
            ]
        }
    }
    const screwedUpCartWithDupes = {
        cart: {
            id: 'dupey-cart-1',
            shopperId: 'shopper-1',
            subtotal: '19.99',
            total: '23.99',
            tax: '4.00',
            taxRate: '20',
            items: [
                {
                    productId: 'p1',
                    quantity: 4,
                    subtotal: '19.99',
                    price: '19.99',
                    title: 'p1'
                },
                {
                    productId: 'p1',
                    quantity: 2,
                    subtotal: '19.99',
                    price: '19.99',
                    title: 'p1'
                }
            ]
        }
    }
    beforeEach(async () => {
        vi.mock('@urql/svelte', async () => {
            let mod = await vi.importActual('@urql/svelte')
            return {
                ...mod,
                subscriptionStore: vi.fn(() => writable(42)),
            }
        })
    })
    afterEach(async () => {
        vi.clearAllMocks()
    })

    describe('given first load', async () => {

        it('should init', async () => {
            let expectCart = initialCart
            let client = createTestClient()
            let store = createCartStore(client)
            use(temporalShop.query('CurrentCart', async (req, res, ctx) => {
                return res.once(ctx.data(expectCart))
            }))
            await new Promise(r => setTimeout(r, 10))
            let actual = get(store)
            expect(subscriptionStore).toHaveBeenCalledOnce()
            expect(subscriptionStore).toHaveBeenCalledWith({
                client: client,
                query: CartDocument,
                variables: {input: {cartId: initialCart.cart.id}}
            })
            expect(actual.data).to.eql(initialCart)
        })

    })
    describe('given initialized cart', async () => {
        let sut: CartStore

        beforeEach(async () => {
            let client = createTestClient()
            sut = createCartStore(client)
            use(temporalShop.query('CurrentCart', async (req, res, ctx) => {
                return res.once(ctx.data(initialCart))
            }))
            await new Promise(r => setTimeout(r, 10))
        })
        it('should add new item', async () => {
            const addItemCommand = {productId: 'p2', quantity: 2}

            let setCartItemsVarbs: SetCartItemsInput | undefined
            use(temporalShop.mutation('SetCartItems',
                (req, res, ctx) => {

                    setCartItemsVarbs = req.variables?.input
                    let postMutation = {
                        ...initialCart.cart, items: [
                            ...initialCart.cart.items, {...addItemCommand, subtotal: '', price: '', title: ''}
                        ]
                    }

                    return res.once(ctx.data({
                        setCartItems: postMutation
                    }))
                }))

            await sut.addItemToCart(addItemCommand)
            await new Promise(r => setTimeout(r, 10))

            expect(setCartItemsVarbs).not.to.be.undefined
            expect(setCartItemsVarbs?.items?.length).to.eq(2)
            let p1 = find(setCartItemsVarbs?.items, {productId: 'p1'})
            let p2 = find(setCartItemsVarbs?.items, {productId: 'p2'})
            expect(p1?.quantity).to.eq(1)
            expect(p2?.quantity).to.eq(addItemCommand.quantity)
        })
        it('should add to existing item', async () => {
            const addItemCommand = {productId: 'p1', quantity: 2}
            let setCartItemsVarbs: SetCartItemsInput | undefined

            use(temporalShop.mutation('SetCartItems',
                (req, res, ctx) => {

                    setCartItemsVarbs = req.variables?.input
                    let postMutation = {
                        ...initialCart.cart, items: [
                            ...initialCart.cart.items, {...addItemCommand, subtotal: '', price: '', title: ''}
                        ]
                    }
                    return res.once(ctx.data({
                        setCartItems: postMutation
                    }))
                }))
            await sut.addItemToCart(addItemCommand)
            await new Promise(r => setTimeout(r, 10))
            expect(setCartItemsVarbs).not.to.be.undefined
            expect(setCartItemsVarbs?.items?.length).to.eq(1)
            let p1 = find(setCartItemsVarbs?.items, {productId: 'p1'})
            expect(p1?.quantity).to.eq(3)
        })
        it('should set new item', async () => {
            const setItemCommand = {productId: 'p2', quantity: 2}
            let setCartItemsVarbs: SetCartItemsInput | undefined

            use(temporalShop.mutation('SetCartItems',
                (req, res, ctx) => {

                    setCartItemsVarbs = req.variables?.input
                    let postMutation = {
                        ...initialCart.cart, items: [
                            ...initialCart.cart.items, {...setItemCommand, subtotal: '', price: '', title: ''}
                        ]
                    }
                    return res.once(ctx.data({
                        setCartItems: postMutation
                    }))
                }))

            await sut.setItemQuantity(setItemCommand)
            await new Promise(r => setTimeout(r, 10))
            expect(setCartItemsVarbs).not.to.be.undefined
            expect(setCartItemsVarbs?.items?.length).to.eq(2)
            let p1 = find(setCartItemsVarbs?.items, {productId: 'p1'})
            let p2 = find(setCartItemsVarbs?.items, {productId: 'p2'})
            expect(p1?.quantity).to.eq(1)
            expect(p2?.quantity).to.eq(setItemCommand.quantity)
        })
        it('should put existing item', async () => {
            const setItemCommand = {productId: 'p1', quantity: 4}
            let setCartItemsVarbs: SetCartItemsInput | undefined

            use(temporalShop.mutation('SetCartItems',
                (req, res, ctx) => {

                    setCartItemsVarbs = req.variables?.input
                    let postMutation = {
                        ...initialCart.cart, items: [
                            ...initialCart.cart.items, {...setItemCommand, subtotal: '', price: '', title: ''}
                        ]
                    }
                    return res.once(ctx.data({
                        setCartItems: postMutation
                    }))
                }))

            await sut.setItemQuantity(setItemCommand)
            await new Promise(r => setTimeout(r, 10))
            expect(setCartItemsVarbs).not.to.be.undefined
            expect(setCartItemsVarbs?.items?.length).to.eq(1)
            let p1 = find(setCartItemsVarbs?.items, {productId: 'p1'})
            expect(p1?.quantity).to.eq(4)
        })
    })
    describe('given cart with duplicate items', async () => {
        let sut: CartStore

        beforeEach(async () => {
            let client = createTestClient()
            sut = createCartStore(client)
            use(temporalShop.query('CurrentCart', async (req, res, ctx) => {
                return res.once(ctx.data(screwedUpCartWithDupes))
            }))
            await new Promise(r => setTimeout(r, 10))
        })
        it('should add new item', async () => {
            const addItemCommand = {productId: 'p2', quantity: 2}
            let setCartItemsVarbs: SetCartItemsInput | undefined
            expect(get(sut).data.cart.id).to.eq(screwedUpCartWithDupes?.cart.id)

            use(temporalShop.mutation('SetCartItems',
                (req, res, ctx) => {

                    setCartItemsVarbs = req.variables?.input
                    let postMutation = {
                        ...screwedUpCartWithDupes.cart, items: [
                            ...screwedUpCartWithDupes.cart.items, {
                                ...addItemCommand,
                                subtotal: '',
                                price: '',
                                title: ''
                            }
                        ]
                    }
                    return res.once(ctx.data({
                        setCartItems: postMutation
                    }))
                }))

            await sut.addItemToCart(addItemCommand)
            await new Promise(r => setTimeout(r, 10))
            expect(setCartItemsVarbs, 'setCartItemVarbs').not.to.be.undefined
            expect(setCartItemsVarbs?.items?.length, 'should have removed duped').to.eq(2)
            let p1 = find(setCartItemsVarbs?.items, {productId: 'p1'})
            let p2 = find(setCartItemsVarbs?.items, {productId: 'p2'})
            expect(p1?.quantity, 'p1 qty').to.eq(4 + 2) // combines current items having same productIds wins
            expect(p2?.quantity, 'p2 qty').to.eq(addItemCommand.quantity)
        })
    })
})