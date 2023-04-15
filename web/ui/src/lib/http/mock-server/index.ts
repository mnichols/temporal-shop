import {graphql, RequestHandler, rest} from "msw";
import {setupServer} from "msw/node";
import {PUBLIC_GRAPHQL_URL} from "$env/static/public";

export const temporalShop = graphql.link(PUBLIC_GRAPHQL_URL)

const handlers = [
    rest.post('http://testing.com/sub', (req, res, ctx) => {
        return res(ctx.json({ok:true}))
    })
]
export const createMockServer = ( ...handlers: Array<RequestHandler>) => {
    return setupServer(...handlers)
}
export const mockServer = createMockServer()
export const use = (...handlers: Array<RequestHandler>) => {
    mockServer.use(...handlers)
}