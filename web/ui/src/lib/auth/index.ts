import {makeOperation} from "@urql/core";
import { authExchange } from '@urql/exchange-auth'
import type {Exchange} from "urql";
const localStorageToken = 'token'
const refreshStorageToken = 'refreshToken'
import { browser } from '$app/environment';
import {apiFetch} from "../http";
import { go } from '$lib/nav'
import { base } from '$app/paths'

export const createAuthExchange = (): Exchange => {
    return authExchange({
        getAuth,
        addAuthToOperation,
        didAuthError,
    })
}
const isLogin = () => {
    return browser && window.location.pathname.includes('login')
}
const getAuth = async ({ authState }) => {

    if (!authState && browser) {
        const token = window.localStorage.getItem(localStorageToken)
        if (token ) {
            return { token, }
        }
    }
    let result = await logout()
    return result
}
const addAuthToOperation = ({ authState, operation }) => {
    console.log('addAuthToOperation', authState, operation)
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
        (error.response && error.response.status === 401)
}
const logout = async () => {
    if (browser) {
        localStorage.removeItem(localStorageToken)
        if(!isLogin()) {
            return await go(`/login`)
        }
    }
    return null
}

interface LoginRequest {
    email: string
    password: string
}
interface LoginResponse {
    email: string
    token: string
}
export const login = async (params: LoginRequest): Promise<void> => {
    if(!browser) {
        console.error('only works in browser')
        return
    }
    let res = await apiFetch({ url: '{scheme}{host}/login' }, {
        method: 'POST',
        body: JSON.stringify(params),
    })
    if (res.response.status === 200) {
        let result : LoginResponse = await res.response.json()
        window.localStorage.setItem(localStorageToken, result.token)
        return go('/')
    } else {
        console.error('failed to login', await res.response.text())
    }

}