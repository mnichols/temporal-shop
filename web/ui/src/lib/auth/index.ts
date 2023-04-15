// const localStorageToken = 'token'
// const refreshStorageToken = 'refreshToken'
// import { browser } from '$app/environment';
// import {apiFetch} from "../http";
// import { go } from '$lib/nav'
//
// export interface Authenticator {
//     getAuth(): Promise<Authentication>;
//     logout(): Promise<void | null | undefined>;
//     login(params: LoginRequest): Promise<void | null | undefined>;
//     assert(): Promise<Authentication | void>;
//     assertResponse(response: Response): Promise<Response | void>
//     withAuth(headers: Headers | Record<string, string> ): (Headers | Record<string, string>)
// }
// const isLogin = () => {
//     return browser && window.location.pathname.includes('login')
// }
// //
// // const didAuthError = ({ error }) => {
// //     return error.graphQLErrors.some(e => e.extensions?.code === 'FORBIDDEN') ||
// //         error.graphQLErrors.some(e => e.extensions?.code === 'UNAUTHORIZED') ||
// //         (error.response && error.response.status === 401)
// // }
//
// export interface Authentication {
//     token?: string
//     ok: boolean
// }
//
// interface LoginRequest {
//     email: string
//     password: string
// }
// interface LoginResponse {
//     email: string
//     token: string
// }
//
// export function createNoOpAuthenticator(): Authenticator {
//     /*
//         {
//       "sub": "ignore",
//       "email": "fake@temporal.io",
//       "iat": 1516239022
//     }
//      */
//     const JWT = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJpZ25vcmUiLCJlbWFpbCI6ImZha2VAdGVtcG9yYWwuaW8iLCJpYXQiOjE1MTYyMzkwMjJ9.F431TlnhKWOPksWBtIZIyPSx9vmyd4RqU3Y9brH-WlE'
//     return {
//         assert(): Promise<Authentication|void> {
//             return Promise.resolve({ok:true});
//         }, getAuth(): Promise<Authentication> {
//             return Promise.resolve({ok:true});
//         }, login(params: LoginRequest): Promise<void | null | undefined> {
//             return Promise.resolve(undefined);
//         }, logout(): Promise<void | null | undefined> {
//             return Promise.resolve(undefined);
//         }, assertResponse(response: Response): Promise<Response | void> {
//             return Promise.resolve(response);
//         }, withAuth(headers: Headers | Record<string, string>): (Headers | Record<string,string>) {
//             return headers
//         }
//     }
// }
//
// export function createDefaultAuthenticator(): Authenticator {
//     const getAuth = async ():Promise<Authentication> => {
//         if (browser) {
//             const token = window.localStorage.getItem(localStorageToken)
//             if (token ) {
//                 return { token,ok:true }
//             }
//         }
//         return {ok:false}
//     }
//
//     const logout = async (): Promise<void> => {
//         if (browser) {
//             localStorage.removeItem(localStorageToken)
//             if(!isLogin()) {
//                 return await go(`/login`)
//             }
//         }
//     }
//
//     const login = async (params: LoginRequest): Promise<void> => {
//         if(!browser) {
//             console.error('only works in browser')
//             return
//         }
//         // TODO get this url from server
//         let res = await apiFetch({ url: '{scheme}{host}/login' }, {
//             method: 'POST',
//             body: JSON.stringify(params),
//         })
//         if (res.response.status === 200) {
//             let result : LoginResponse = await res.response.json()
//             window.localStorage.setItem(localStorageToken, result.token)
//             return go('/')
//         } else {
//             console.error('failed to login', await res.response.text())
//         }
//
//     }
//     const assert = async (): Promise<Authentication | void> => {
//         let auth = await getAuth()
//         if (!auth?.token) {
//             console.log('token not found', 'logging out')
//             return await logout()
//         }
//         return { token: auth?.token, ok: !!auth?.token}
//     }
//     const assertResponse = async(response: Response): Promise<Response | void> => {
//         console.log('assertResponse', response.status)
//         if(response.status === 401) {
//             return logout()
//         }
//         return response
//     }
//     return {
//         getAuth,
//         login,
//         logout,
//         assert,
//         assertResponse,
//         withAuth(headers: Headers | Record<string, string>): Headers | Record<string, string> {
//
//         }
//     }
// }