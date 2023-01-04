// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type Cart struct {
	ID string `json:"id"`
}

type CartItem struct {
	ProductID string `json:"productId"`
}

type Game struct {
	ID       string `json:"id"`
	Product  string `json:"product"`
	Category string `json:"category"`
	ImageURL string `json:"imageUrl"`
	Price    string `json:"price"`
}

type Inventory struct {
	Games []*Game `json:"games"`
}

type InventoryInput struct {
	CategoryID *string `json:"categoryId"`
}

type Shopper struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	InventoryID string `json:"inventoryId"`
}

type ShopperInput struct {
	ShopperID *string `json:"shopperId"`
}
