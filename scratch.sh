# logging in (getting token)
export JWT = $(http --verify=no https://localhost:8080/api/login email=mike.nichols@temporal.io | jq -r .token)

# subscribing
cat sub.json| http -v --verify=no https://localhost:8080/sub Authorization:"Bearer $JWT" Accept:"text/event-stream" Content-Type:"application/json"

# publishing with Temporal signal to cart workflow
tctctl workflow signal -w "cart_9e33e4c43a46b38fb9554226d473c64aee704b2d102b58ff28b3bc67bcfbb91cc56c72d630cd75fb6a80" --name "temporal_shop.commands.v1.SetCartItemsRequest" --input-file set_cart_items.json
