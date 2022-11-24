import { browser } from '$app/environment'
import { noop } from 'svelte/internal'

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
    params?: {}
    request?: typeof fetch
    onError?: ErrorCallback,
    isBrowser?: boolean,
}
type APIResponse = {
    response: Response,
    data: unknown,
}
export const apiFetch = async <T>(
    url: string,
    opts: RequestOpts
): Promise<APIResponse> => {
    const {
        params = {},
        request = fetch,
        onError = noop(),
        isBrowser = browser,
    } = opts
    let requestOpts = { }
    let res = await request(url, requestOpts)
    return {
        response: res,
        data: {},
    }
}
// import { handleError as handleRequestError } from '../errors'
// import { isFunction } from '../utilities/is-function'
// // import { toURL } from './to-url'
//

//
// export type RetryCallback = (retriesRemaining: number) => void
//

//
//
// type toURLParams = Parameters<typeof toURL>
//
// type RequestFromAPIOptions = {
//   params?: toURLParams[1]
//   request?: typeof fetch
//   options?: Parameters<typeof fetch>[1]
//   token?: string
//   onRetry?: RetryCallback
//   onError?: ErrorCallback
//   notifyOnError?: boolean
//   handleError?: typeof handleRequestError
//   shouldRetry?: boolean
//   retryInterval?: number
//   isBrowser?: boolean
// }
//
// export const isAPIError = (obj: unknown): obj is APIError =>
//   (obj as APIError)?.message !== undefined &&
//   typeof (obj as APIError)?.message === 'string'
//
// /**
//  *  A utility method for making requests to the Temporal API.
//  *
//  * @param endpoint The path of the API endpoint you want to request data from.
//  *
//  * @param options.params Query (or search) parameters to be suffixed to the
//  * path.
//  * @param options.token Shorthand for a `nextPageToken` query parameter.
//  * @param options.request A replacement for the native `fetch` function.
//  *
//  * @returns Promise with the response from the API parsed into an object.
//  */
// export const apiFetch = async <T>(
//   endpoint: toURLParams[0],
//   init: RequestFromAPIOptions = {},
//   retryCount = 10,
// ): Promise<T> => {
//   const {
//     params = {},
//     request = fetch,
//     token,
//     shouldRetry = false,
//     notifyOnError = true,
//     handleError = handleRequestError,
//     onRetry = noop,
//     onError,
//     retryInterval = 5000,
//     isBrowser = browser,
//   } = init
//   let { options } = init
//
//   const nextPageToken = token ? { next_page_token: token } : {}
//   const query = new URLSearchParams({
//     ...params,
//     ...nextPageToken,
//   })
//   const url = toURL(endpoint, query)
//
//   try {
//     options = withSecurityOptions(options, isBrowser)
//     //options = await withAuth(options, isBrowser);
//
//     const response = await request(url, options)
//     const body = await response.json()
//
//     const { status, statusText } = response
//
//     if (!response.ok) {
//       if (onError && isFunction(onError)) {
//         onError({ status, statusText, body })
//       } else {
//         throw {
//           statusCode: response.status,
//           statusText: response.statusText,
//           response,
//           message: body?.message ?? response.statusText,
//         } as NetworkError
//       }
//     }
//
//     return body
//   } catch (error: unknown) {
//     if (notifyOnError) {
//       handleError(error)
//
//       if (shouldRetry && retryCount > 0) {
//         return new Promise((resolve) => {
//           const retriesRemaining = retryCount - 1
//           onRetry(retriesRemaining)
//           setTimeout(() => {
//             resolve(apiFetch(endpoint, init, retriesRemaining))
//           }, retryInterval)
//         })
//       }
//       throw Error('foo')
//     } else {
//       throw error
//     }
//   }
// }
//

const withSecurityOptions = (
    options: RequestInit,
    isBrowser: browser,
): RequestInit => {
    const opts: RequestInit = { credentials: 'include', ...options }
    opts.headers = withCsrf(options?.headers, isBrowser)
    return  { credentials: 'include', ...options}

}
const withCsrf = (headers?: HeadersInit, isBrowser: boolean = browser): HeadersInit => {
    let h = new Headers(headers)
    const csrfCookie = '_csrf='
    const csrfHeader = 'X-CSRF-TOKEN'
    if (!isBrowser) {
        return h
    }
    if(h.has(csrfHeader)) {
        return h
    }
    try {
        const cookies = document.cookie.split(';')
        let token = cookies.find((c) => c.includes(csrfCookie))
        if (token) {
            token = token.trim().slice(csrfCookie.length)
            h.set(csrfHeader, token)
        }
    } catch (err) {
        console.error(err)
    }
    return h
}
