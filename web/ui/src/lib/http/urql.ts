import { browser } from '$app/environment'
import { createClient, setContextClient,dedupExchange, cacheExchange, fetchExchange } from '@urql/svelte'
import { createAuthExchange } from '../auth'
import type { Client } from '@urql/svelte'
import { PUBLIC_GRAPHQL_URL} from '$env/static/public'
import {withSecurityOptions} from './api-fetch.js'

const authExchange = createAuthExchange()
export const createGraphQLClient = (): Client => {
    return createClient({
        url: PUBLIC_GRAPHQL_URL,
        fetchOptions: (): RequestInit => {
            let i = withSecurityOptions({}, browser)
            i.headers = Object.entries(new Headers(i.headers).entries())
            return i
        },
        exchanges: [
            dedupExchange,
            cacheExchange,
            authExchange,
            fetchExchange,
        ],
    })
}

