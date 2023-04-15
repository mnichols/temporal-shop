import {cacheExchange} from "@urql/exchange-graphcache";

export const createCacheExchange = () => {
    let cache = cacheExchange({
        // specify custom keys
        // https://formidable.com/open-source/urql/docs/graphcache/normalized-caching/#custom-keys-and-non-keyable-entities
        keys: {
            CartItem: data => data.productId,
            User: data => data.email,
            Inventory: data => null,
        },
    })
    return cache
}
