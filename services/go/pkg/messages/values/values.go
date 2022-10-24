package values

type HTTPResponse struct {
	StatusCode int `json:"status_code"`
}
type Product struct {
	ProductID string `json:"product_id"`
	Name      string `json:"name"`
	Category  string `json:"category"`
	ImageURL  string `json:"image_url"`
	Price     int    `json:"price"`
}
type Promotion struct {
	DiscountPercentage int    `json:"discount_percentage" validate:"gte=0"`
	Code               string `json:"code"`
}
type Purchase struct {
	Product  *Product `json:"product" validate:"required"`
	Quantity int      `json:"quantity" validate:"gte=0"`
}
