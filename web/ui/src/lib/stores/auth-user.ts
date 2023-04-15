/*
auth-user
babysits the presence of token at localStorage('token')

 */

import {
    get,
    writable,
    derived, asyncDerived, asyncReadable,
} from "@square/svelte-store";
import type { Writable} from '@square/svelte-store'
import {go} from "$lib/nav";
import {browser} from "$app/environment";
import {apiFetch} from "../http";
const tokenName = 'token'
import { persisted } from 'svelte-local-storage-store'
import {Client, getContextClient, gql,  queryStore} from "@urql/svelte";
import type {UserInput, User, QueryUserArgs, CurrentUserQuery, CurrentUserQueryVariables} from "../../gql";
import {CurrentUserDocument} from "$gql";
import {getContext, setContext} from "svelte";
import type { OperationResultStore, Pausable, QueryArgs,} from '@urql/svelte'

const _contextKey = '$$_user'

interface LoginRequest {
    email: string
    password: string
}
interface LoginResponse {
    email: string
    token: string
}

const isLogin = () => {
    return browser && window.location.pathname.includes('login')
}

// internal tokenStore for managing localStorage token between sessions
let _tokenStore: Writable<string | null | undefined> = writable()

// response storage from /login page
const loginResponse = writable<LoginResponse | null>()

if(browser) {
    _tokenStore = persisted(tokenName,undefined)
}
_tokenStore.subscribe(t => { console.log('_tokenStore#persisted','write token', t)})
const sanitizeToken = (token?: string | null) : string | undefined => {
    if(!token) {
        return undefined
    }
    let replaced = token.replaceAll('Bearer ', '').trim()
    return replaced
}
loginResponse.subscribe(res => {
    if(res?.token) {
        let t = sanitizeToken(res?.token)
        console.log('_tokenStore.set', t)
        _tokenStore.set(t)
    }
})

/** stores **/
type UserStore = OperationResultStore<CurrentUserQuery, CurrentUserQueryVariables | undefined> & Pausable
export const createUserStore = (client: Client, userInput?: UserInput):UserStore => {
    let args: QueryArgs<CurrentUserQuery, CurrentUserQueryVariables> = {
        client: client,
        query:CurrentUserDocument,
    }
    if(userInput) {
        args.variables = { input: userInput }
    }
    let u = queryStore(args)
    u.subscribe(arg => {
        if(arg.fetching) {
            return
        }

        if(arg?.data?.user) {
            _tokenStore.set(sanitizeToken(arg?.data?.user?.token))
        }
    })

    return queryStore(args)
}

export const doLogin = async (params: LoginRequest) => {
    if(!browser) {
        console.error('only works in browser')
        return
    }
    // TODO get this url from server
    let res = await apiFetch({ url: '{scheme}{host}/login' }, {
        method: 'POST',
        body: JSON.stringify(params),
    })
    if(res.response.status > 299) {
        console.error('failed to login', res.response.status, await res?.response?.text())
    }
    let result : LoginResponse = await res.response.json()
    loginResponse.set(result)
    return await go('/')
}
export const goLogin = async () : Promise<void> => {
    if(isLogin()) {
        return
    }
    console.log('#goLogin')
    return await go('/login')
}
export const doLogout = async () => {
    // return await storedToken.clear()
}

export function withAuth (headers: HeadersInit | Headers | Record<string, string>): Headers  {
    let h = new Headers(headers)
    if(h.has('authorization')) {
        console.log('authorization already here')
        return h
    }
    h.set('authorization', `Bearer ${get(_tokenStore)}`)
    return h
}
export const getContextUser = (): UserStore => {
    const out = getContext(_contextKey);
    if (process.env.NODE_ENV !== 'production' && !out) {
        throw new Error(
            'No Cart was found in Svelte context. Did you forget to call setContextCart?'
        );
    }

    return out as UserStore;
}
export const setContextUser = (user: UserStore): void => {
    setContext(_contextKey, user)
}
export const tokenStore = _tokenStore
