/** @type {import('./$types').PageLoad} */
export function load({ params }) {
    return {
        categoryId: params.category,
    }
}