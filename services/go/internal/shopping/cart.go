package shopping

import (
	"context"
	"fmt"
	commands "github.com/temporalio/temporal-shop/services/go/api/generated/temporal_shop/commands/v1"
	queries "github.com/temporalio/temporal-shop/services/go/api/generated/temporal_shop/queries/v1"
	values "github.com/temporalio/temporal-shop/services/go/api/generated/temporal_shop/values/v1"
	"google.golang.org/protobuf/proto"
)

type ShoppingCart struct {
	commands   []proto.Message
	taxRateBPS int64
	args       ShoppingCartArgs
}

type ShoppingCartArgs struct {
	ID         string
	ShopperID  string
	TaxRateBPS int64
}

type FetchGames func(productIDs []string) ([]*values.Game, error)

func NewShoppingCart(args ShoppingCartArgs) *ShoppingCart {
	return &ShoppingCart{
		args:       args,
		commands:   make([]proto.Message, 0),
		taxRateBPS: 0,
	}
}
func (c *ShoppingCart) Empty() (*queries.GetCartResponse, error) {
	state := &queries.GetCartResponse{
		CartId:              c.args.ID,
		ShopperId:           c.args.ShopperID,
		SubtotalCents:       0,
		TaxCents:            0,
		TaxRateBps:          0,
		TotalCents:          0,
		ProductIdToQuantity: make(map[string]int64),
		ProductIdToGame:     make(map[string]*values.Game),
	}
	return state, nil
}
func (c *ShoppingCart) Calculate(gamesFetcher FetchGames) (*queries.GetCartResponse, error) {
	state, err := c.Empty()
	if err != nil {
		return nil, err
	}
	if gamesFetcher == nil {
		return nil, fmt.Errorf("gamesFetcher is required")
	}
	out, err := c.sync(state, gamesFetcher)
	if err != nil {
		return nil, err
	}
	return out, err
}

func (c *ShoppingCart) Append(msg ...proto.Message) {
	// only allow one message for now
	c.commands = append([]proto.Message{}, msg...)
	//c.commands = append(c.commands, msg...)
}

func (c *ShoppingCart) sync(state *queries.GetCartResponse, fetcher FetchGames) (*queries.GetCartResponse, error) {

	for _, cmd := range c.commands {
		switch v := cmd.(type) {
		case *commands.SetCartItemsRequest:
			if err := c.setCartItems(state, fetcher, v); err != nil {
				return state, err
			}
		default:
			return state, fmt.Errorf("unknown command type %T", v)
		}
	}
	return c.calculate(state)
}
func (c *ShoppingCart) setCartItems(state *queries.GetCartResponse, fetcher FetchGames, cmd *commands.SetCartItemsRequest) error {
	var productIdsToFetch []string
	var includeMap = make(map[string]struct{})
	state.ProductIdToQuantity = make(map[string]int64)
	for pid, qty := range cmd.ProductIdsToQuantity {
		if qty > 0 {
			state.ProductIdToQuantity[pid] = qty
		}
	}
	gameMap := map[string]*values.Game{}

	for pid := range cmd.ProductIdsToQuantity {
		if _, exists := state.ProductIdToGame[pid]; !exists && cmd.ProductIdsToQuantity[pid] > 0 {
			productIdsToFetch = append(productIdsToFetch, pid)
		}
		includeMap[pid] = struct{}{}
	}
	// copy over from state cache
	for _, g := range state.ProductIdToGame {
		if _, exists := includeMap[g.Id]; exists {
			gameMap[g.Id] = g
		}
	}
	if len(productIdsToFetch) > 0 {
		var games []*values.Game
		var err error
		games, err = fetcher(productIdsToFetch)
		if err != nil {
			return err
		}
		for _, g := range games {
			if cmd.ProductIdsToQuantity[g.Id] > 0 {
				gameMap[g.Id] = g
			}
		}
	}

	state.ProductIdToGame = gameMap

	return nil
}

func (c *ShoppingCart) calculate(next *queries.GetCartResponse) (*queries.GetCartResponse, error) {
	next.TotalCents = 0
	next.SubtotalCents = 0
	if next.TaxRateBps == 0 {
		next.TaxRateBps = DefaultTaxRateBPS
	}
	for k, g := range next.ProductIdToGame {
		next.SubtotalCents = next.SubtotalCents + (g.PriceCents * next.ProductIdToQuantity[k])
	}
	next.TotalCents = CalculateTotalCents(next.SubtotalCents, next.TaxRateBps)
	next.TaxCents = next.TotalCents - next.SubtotalCents
	return next, nil
}

func (c *ShoppingCart) Calculate2(ctx context.Context, cmd *commands.CalculateShoppingCartRequest) (*queries.GetCartResponse, error) {
	state, err := c.Empty()
	if err != nil {
		return nil, err
	}

	for gid, qty := range cmd.ProductIdsToQuantity {
		if qty > 0 {
			g, exists := cmd.ProductIdToGame[gid]
			if !exists {
				return nil, fmt.Errorf("failed to find game with product id [%s] in collection", gid)
			}
			state.ProductIdToGame[g.Id] = g
			state.ProductIdToQuantity[g.Id] = qty
		}
	}
	return c.calculate(state)
}
