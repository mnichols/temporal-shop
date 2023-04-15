import type { Handle } from '@sveltejs/kit';
import {go} from "./lib/nav";

export const handle: Handle = async ({ event, resolve }) => {
    const response = await resolve(event, {});
    console.log('handleserverhook', response.status)
    // if(response.status == 401) {
    //     return await go('/login')
    // }
    return response;
};