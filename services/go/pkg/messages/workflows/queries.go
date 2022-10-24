package workflows

const QueryOrder = "query_order"

type QueryOrderResult struct {
	PaymentIntentID string
	OrderID         string
}
