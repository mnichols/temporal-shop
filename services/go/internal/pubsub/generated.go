// Code generated by github.com/Khan/genqlient, DO NOT EDIT.

package pubsub

import (
	"context"

	"github.com/Khan/genqlient/graphql"
)

type PublishCartInput struct {
	Id            string                 `json:"id"`
	ShopperId     string                 `json:"shopperId"`
	Items         []PublishCartItemInput `json:"items"`
	SubtotalCents int                    `json:"subtotalCents"`
	TaxRateBps    int                    `json:"taxRateBps"`
	TotalCents    int                    `json:"totalCents"`
	TaxCents      int                    `json:"taxCents"`
}

// GetId returns PublishCartInput.Id, and is useful for accessing the field via an interface.
func (v *PublishCartInput) GetId() string { return v.Id }

// GetShopperId returns PublishCartInput.ShopperId, and is useful for accessing the field via an interface.
func (v *PublishCartInput) GetShopperId() string { return v.ShopperId }

// GetItems returns PublishCartInput.Items, and is useful for accessing the field via an interface.
func (v *PublishCartInput) GetItems() []PublishCartItemInput { return v.Items }

// GetSubtotalCents returns PublishCartInput.SubtotalCents, and is useful for accessing the field via an interface.
func (v *PublishCartInput) GetSubtotalCents() int { return v.SubtotalCents }

// GetTaxRateBps returns PublishCartInput.TaxRateBps, and is useful for accessing the field via an interface.
func (v *PublishCartInput) GetTaxRateBps() int { return v.TaxRateBps }

// GetTotalCents returns PublishCartInput.TotalCents, and is useful for accessing the field via an interface.
func (v *PublishCartInput) GetTotalCents() int { return v.TotalCents }

// GetTaxCents returns PublishCartInput.TaxCents, and is useful for accessing the field via an interface.
func (v *PublishCartInput) GetTaxCents() int { return v.TaxCents }

type PublishCartItemInput struct {
	ProductId     string `json:"productId"`
	Quantity      int    `json:"quantity"`
	SubtotalCents int    `json:"subtotalCents"`
	PriceCents    int    `json:"priceCents"`
	Title         string `json:"title"`
}

// GetProductId returns PublishCartItemInput.ProductId, and is useful for accessing the field via an interface.
func (v *PublishCartItemInput) GetProductId() string { return v.ProductId }

// GetQuantity returns PublishCartItemInput.Quantity, and is useful for accessing the field via an interface.
func (v *PublishCartItemInput) GetQuantity() int { return v.Quantity }

// GetSubtotalCents returns PublishCartItemInput.SubtotalCents, and is useful for accessing the field via an interface.
func (v *PublishCartItemInput) GetSubtotalCents() int { return v.SubtotalCents }

// GetPriceCents returns PublishCartItemInput.PriceCents, and is useful for accessing the field via an interface.
func (v *PublishCartItemInput) GetPriceCents() int { return v.PriceCents }

// GetTitle returns PublishCartItemInput.Title, and is useful for accessing the field via an interface.
func (v *PublishCartItemInput) GetTitle() string { return v.Title }

// PublishCartPublishCart includes the requested fields of the GraphQL type Cart.
type PublishCartPublishCart struct {
	Id string `json:"id"`
}

// GetId returns PublishCartPublishCart.Id, and is useful for accessing the field via an interface.
func (v *PublishCartPublishCart) GetId() string { return v.Id }

// PublishCartResponse is returned by PublishCart on success.
type PublishCartResponse struct {
	PublishCart PublishCartPublishCart `json:"publishCart"`
}

// GetPublishCart returns PublishCartResponse.PublishCart, and is useful for accessing the field via an interface.
func (v *PublishCartResponse) GetPublishCart() PublishCartPublishCart { return v.PublishCart }

// __PublishCartInput is used internally by genqlient
type __PublishCartInput struct {
	Input PublishCartInput `json:"input"`
}

// GetInput returns __PublishCartInput.Input, and is useful for accessing the field via an interface.
func (v *__PublishCartInput) GetInput() PublishCartInput { return v.Input }

func PublishCart(
	ctx context.Context,
	client graphql.Client,
	input PublishCartInput,
) (*PublishCartResponse, error) {
	req := &graphql.Request{
		OpName: "PublishCart",
		Query: `
mutation PublishCart ($input: PublishCartInput!) {
	publishCart(input: $input) {
		id
	}
}
`,
		Variables: &__PublishCartInput{
			Input: input,
		},
	}
	var err error

	var data PublishCartResponse
	resp := &graphql.Response{Data: &data}

	err = client.MakeRequest(
		ctx,
		req,
		resp,
	)

	return &data, err
}
