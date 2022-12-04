import {makeOperation} from "@urql/core";
import { authExchange } from '@urql/exchange-auth'
import type {Exchange} from "urql";
const localStorageToken = 'token'
const refreshStorageToken = 'refreshToken'
import { browser } from '$app/environment';


export const createAuthExchange = (): Exchange => {
    return authExchange({
        getAuth,
        addAuthToOperation,
        didAuthError,
    })
}
const getAuth = async ({ authState }) => {
    if (!authState && browser) {
        const token = window.localStorage.getItem(localStorageToken)
        if (token ) {
            return { token, }
        }
        return null
    }

    logout()
    return null
}
const addAuthToOperation = ({ authState, operation }) => {
    if (!authState || !authState.token) {
        return operation
    }

    const fetchOptions =
        typeof operation.context.fetchOptions === 'function'
            ? operation.context.fetchOptions()
            : operation.context.fetchOptions || {}

    return makeOperation(operation.kind, operation, {
        ...operation.context,
        fetchOptions: {
            ...fetchOptions,
            headers: {
                ...fetchOptions.headers,
                Authorization: authState.token,
            },
        },
    })
}

const didAuthError = ({ error }) => {
    return error.graphQLErrors.some(e => e.extensions?.code === 'FORBIDDEN') ||
        error.graphQLErrors.some(e => e.extensions?.code === 'UNAUTHORIZED') ||
        error.graphQLErrors.some(
            e => e.response.status === 401,
        )
}
const logout = () => {
    if (browser) {
        localStorage.removeItem(localStorageToken)
    }
}