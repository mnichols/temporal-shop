import type {Exchange, SubscriptionExchangeOpts} from "@urql/svelte";
import {subscriptionExchange} from "@urql/core";
import {fetchSSE} from "../fetch-sse";
import {EventStreamContentType} from "@microsoft/fetch-event-source";
import {FatalError, RetriableError} from "$lib/http/errors";
import { Logger } from '$log'

export const createSubscriptionExchange = (mustUseSubscriptionsURL?: string, opts?: SubscriptionExchangeOpts): Exchange => {
    // const log = console.debug.bind(console,'DEBUG: subscriptions-client')
    return subscriptionExchange({
        enableAllOperations: false,
        forwardSubscription: (fetchBody, operation) => {
            Logger.debug( fetchBody, 'forwardSubscription')
            return {
                subscribe: sink => {
                    let controller = new AbortController()
                    const fetchOptions =
                        typeof operation.context.fetchOptions === 'function'
                            ? operation.context.fetchOptions()
                            : operation.context.fetchOptions ?? {}
                    const headers = Object.entries(fetchOptions.headers ?? {}).reduce(
                        (headers, [key, value]) => ({ ...headers, [key]: value }),
                        {} as Record<string, string>
                    )
                    headers['accept'] = 'text/event-stream'
                    const body = {
                        query: fetchBody.query,
                        variables: fetchBody.variables
                    }
                    const url = mustUseSubscriptionsURL || operation.context.url
                    const err = console.error.bind(console,'fetchSSE')
                    Logger.debug(body, 'initiating fetchSSE %s', url)

                    fetchSSE({url}, {
                        body: JSON.stringify(body),
                        method: 'POST',
                        headers: headers,
                        openWhenHidden: true,
                        onmessage: (ev) => {
                            try {
                                Logger.debug(ev, '#fetch-sse.onmessage: %s','received message')
                                if(ev?.event === 'complete') {
                                    Logger.debug(ev, '#fetch-sse.onmessage: %s','calling sink.complete')
                                    sink.complete()
                                } else if (ev?.data) {
                                    let tmp = JSON.parse(ev.data)
                                    Logger.debug(tmp,'#fetch-sse.onmessage: %s','calling sink.next')
                                    sink.next(tmp)
                                }
                            } catch(e) {
                                Logger.error(e,'#fetch-sse.onmessage: %s','failed to handle `next`')
                            }
                        },
                        signal: controller.signal,
                        onopen: async (response) => {
                            Logger.debug('#fetch-sse.onopen: %s', response.status)
                            if (response.ok && response.headers.get('content-type') === EventStreamContentType) {
                                return Promise.resolve(); // everything's good
                            } else if (response.status >= 400 && response.status < 500 && response.status !== 429) {
                                Logger.error(response, 'onopen failed')
                                // client-side errors are usually non-retriable:
                                throw new FatalError();
                            } else {
                                Logger.error(response,'#fetch-sse.onopen{else}')
                                throw new RetriableError();
                            }
                        },
                        onclose() {
                            Logger.info('#fetch-sse.onclose')
                            // if the server closes the connection unexpectedly, retry:
                            throw new RetriableError();
                        },
                        onerror(err) {
                            Logger.error(err, '#fetch-sse.onerror')
                            if (err instanceof FatalError) {
                                sink.error(err)
                                throw err; // rethrow to stop the operation
                            } else {
                                // do nothing to automatically retry. You can also
                                // return a specific retry interval here.
                            }
                        }
                    }).catch(arg => {
                        Logger.error(arg, '#fetch-sse.catch.failure')
                        sink.complete()
                        // throw new FatalError()
                    })
                    return {
                        unsubscribe: () => {
                            // console.trace()
                            Logger.debug('#fetch-sse#unsubscribe %s', 'calling controller.abort')

                            controller.abort('disconnected sse')
                        }
                    }
                },
            }
        }
    })
}
