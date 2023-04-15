import {go} from "$lib/nav";

export const prerender = true
export const ssr = false

// here is where we will do our auth store wire up to watch for
// changes on user state (eg logged out or 401 returned)
import type { LayoutLoad} from './$types'

export const load: LayoutLoad = async function() {

}