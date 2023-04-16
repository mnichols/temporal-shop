import { createThing } from './testable'
import {createTestClient} from "../http/urql";
import {writable} from "svelte/store";
import {afterEach} from "vitest";
describe('when mocking', async () => {
    afterEach(async () => {
        vi.clearAllMocks()
    })
    it('should call mocked fn', async () => {
        vi.mock('@urql/svelte', async () => {
            let mod = await vi.importActual('@urql/svelte')
            return {
                ...mod,
                subscriptionStore: vi.fn(()=>writable(42)),
            }
        })
        let client = createTestClient()
        let actual = createThing(client)
        actual.subscribe(arg => console.log('RECEIVED SUBSCRIBED VALUE', arg))
        actual.set('DO IT NOW')
    })
})