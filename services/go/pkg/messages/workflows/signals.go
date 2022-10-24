package workflows

const SignalRecomputeRepurchaseSpecification = "repurchase_recompute_repurchase_specification"
const SignalRecompute = "repurchase_customer_purchased"

type RepurchaseSpecificationChanged struct{}
type RepurchaseCustomerPurchased struct{}
