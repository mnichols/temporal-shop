import { browser } from '$app/environment'
import { noop } from 'svelte/internal'
import type { RequireAtLeastOne} from "type-fest"
import parser from 'uri-template'

const HEADER_CONTENT_TYPE='content-type'
export { parse } from 'uri-template'
export const csrfCookie = '_csrf'
export const csrfHeader = 'X-CSRF-TOKEN'

interface Expandable {
    expand(values: Record<string, unknown>): string
}

type TemplatableURL = {
    url?: string
    tpl?: Expandable
    params?: any
}

type RequiredTemplatableURL = RequireAtLeastOne<TemplatableURL, 'url' | 'tpl'>

export type APIError = {
  code: number
  message: string
  details: unknown[]
}
export type APIErrorResponse = {
  status: number
  statusText: string
  body: APIError
}
export type ErrorCallback = (error: APIErrorResponse) => void

type RequestOpts = {
    body?: BodyInit
    method?: string
    headers?: Headers
    request?: typeof fetch
    onError?: ErrorCallback,
    isBrowser?: boolean,
}
type APIResponse = {
    response: Response,
}
export const apiFetch = async <T>(
    url: RequiredTemplatableURL,
    opts: RequestOpts
): Promise<APIResponse> => {
    let actualURL = ''
    if(url.url) {
        let tmp = parser.parse(url.url)
        actualURL = tmp.expand(url.params || {})
    } else if (url.tpl) {
        actualURL = url.tpl.expand(url.params || {})
    }
    const {
        request = fetch,
        onError = noop(),
        isBrowser = browser,
        //headers,
        method = 'GET',
        body,
    } = opts
    let headers = new Headers(opts.headers || {})

    if (!headers.has(HEADER_CONTENT_TYPE)) {
        headers.set(HEADER_CONTENT_TYPE, 'application/json')
    }
  let requestOpts: RequestInit = {
        body,
        method,
        headers,
    }

    requestOpts = withSecurityOptions(requestOpts, browser)
    requestOpts.headers = Object.fromEntries(new Headers(requestOpts.headers).entries())
    console.log('requestOpts', requestOpts.headers)
    let res = await request(actualURL, requestOpts)
    return {
        response: res,
    }
}

export const withSecurityOptions = (
    options: RequestInit = {},
    isBrowser: browser,
): RequestInit => {
    options['credentials'] = 'include'
    options.headers = withCsrf(options?.headers || {}, isBrowser)
    return options
}

export const withCsrf = (headers: HeadersInit, isBrowser: boolean = browser): HeadersInit => {
    if (!isBrowser) {
        return headers || {}
    }
    const h = new Headers(headers)
    if(h.has(csrfHeader)) {
        return h
    }

    try {
        const cookies = document.cookie.split(';')
        let token = cookies.find((c) => c.includes(csrfCookie))
        if(!token) {
            return h
        }
        token = token.trim().slice(csrfCookie.length + 1)
        h.set(csrfHeader, token)
    } catch (err) {
        console.error(err)
    }

    return h
}