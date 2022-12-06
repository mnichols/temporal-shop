import { browser } from '$app/environment'
import { noop } from 'svelte/internal'
import type { RequireAtLeastOne} from "type-fest"
import parser from 'uri-template'
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
    body?: {}
    method?: string
    headers?: HeadersInit
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
        body = {},
        request = fetch,
        onError = noop(),
        isBrowser = browser,
    } = opts
    let requestOpts = { }
    requestOpts = withSecurityOptions(requestOpts, browser)
    let res = await request(actualURL, requestOpts)
    return {
        response: res,
    }
}

export const withSecurityOptions = (
    options: RequestInit,
    isBrowser: browser,
): RequestInit => {
    const opts: RequestInit = { credentials: 'include', ...options }
    opts.headers = withCsrf(options?.headers, isBrowser)
    return opts
}

export const withCsrf = (headers: HeadersInit, isBrowser: boolean = browser): HeadersInit => {
    if (!isBrowser) {

        return headers || {}
    }

    let h = new Headers(headers)
    if(h.has(csrfHeader)) {
        return {}
    }

    try {
        const cookies = document.cookie.split(';')
        let token = cookies.find((c) => c.includes(csrfCookie))
        if(!token) {
            return {}
        }
        token = token.trim().slice(csrfCookie.length + 1)
        h.set(csrfHeader, token)
    } catch (err) {
        console.error(err)
    }

    let out: HeadersInit = {}
    for(const he of h.entries()) {
        out[he[0]] = he[1]
    }
    return out
}