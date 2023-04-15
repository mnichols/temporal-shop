import {temporalShop, use} from "./index";
import { parseISO } from 'date-fns'
import {createClient, queryStore} from "@urql/svelte";
import {createTestClient} from "../urql/index";
import {PingTestDocument} from "$gql"
import type { PingInput, } from '$gql'
import type {Client} from "@urql/core";
import {beforeAll} from "vitest";
describe('given mock server', async () => {
    let client: Client
    beforeAll(async () => {
        client = createTestClient()
    })
    it('should init', async () => {

        let receivedVariables:(PingInput | undefined) = undefined
        use(temporalShop.query('PingTest', async (req, res, ctx) => {
            receivedVariables = req.variables
            const {input } = req.variables
            return res.once(ctx.data({ping: { value: input?.value + ' to you', timestamp: parseISO(input?.timestamp)}}))
        }))
        let input: PingInput = { value: 'hello', timestamp: new Date(1970, 3, 5)}

        let responses = []
        let sut = queryStore({
            client,
            query: PingTestDocument,
            variables: {input},
            requestPolicy: 'network-only',
        })
        sut.subscribe(arg => responses.push(arg))
        // first response is waiting for server reply, so fetching is true
        expect(responses).length(1)
        expect(responses[0].fetching).to.be.true
        await new Promise(r => setTimeout(r, 1))

        // after tick we get the value from the server so fetching is false and we have data
        expect(receivedVariables?.input?.value).to.eq(input?.value)
        expect(receivedVariables?.input?.timestamp).to.include('1970-04-05')
        expect(responses).length(2)
        expect(responses[1].fetching).to.be.false
        expect(responses[1]?.data?.ping?.value).to.eq(`${input?.value} to you`)
    })
})