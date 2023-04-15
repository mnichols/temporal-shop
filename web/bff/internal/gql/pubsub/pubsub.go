package pubsub

import (
	"context"
	"github.com/temporalio/temporal-shop/services/go/pkg/instrumentation/log"
	"github.com/temporalio/temporal-shop/web/bff/internal/gql/graph/model"
	"sync"
)

func NewPubSub() *PubSub {
	result := &PubSub{}
	result.carts = make(map[string]map[string]chan *model.Cart)
	return result
}

type PubSub struct {
	mu        sync.RWMutex
	carts     map[string]map[string]chan *model.Cart
	taskQueue string
}

func (p *PubSub) getCart(cartId, topic string) chan *model.Cart {
	c, exists := p.carts[cartId]
	if !exists {
		return nil
	}
	if ch, exist := c[topic]; exist {
		return ch
	}

	return nil
}
func (p *PubSub) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, c := range p.carts {
		for _, ch := range c {
			close(ch)
		}
	}
	return nil
}
func (p *PubSub) PublishCart(ctx context.Context, cart *model.Cart) error {
	p.mu.RLock()
	defer p.mu.RUnlock()
	logger := log.GetLogger(ctx)
	logger.Info("subscription count", log.Fields{"count": len(p.carts)})
	c := p.carts[cart.ID]
	if c == nil {
		logger.Info("no subscriptions", log.Fields{"cart_id": cart.ID})
		return nil
	}
	logger.Info("publishing cart", log.Fields{"cart_id": cart.ID})
	for topic, ch := range c {
		select {
		case ch <- cart:
		//sent
		default:
			close(ch)
			delete(c, topic)
		}
	}
	return nil
}
func (p *PubSub) SubscribeCart(ctx context.Context, topic, cartId string) chan *model.Cart {
	p.mu.Lock()
	defer p.mu.Unlock()
	logger := log.GetLogger(ctx)
	logger.Info("subscribing", log.Fields{"topic": topic, "cart_id": cartId})

	var ch chan *model.Cart
	if ch = p.getCart(cartId, topic); ch == nil {
		if _, exists := p.carts[cartId]; !exists {
			p.carts[cartId] = make(map[string]chan *model.Cart)
		}
		ch = make(chan *model.Cart)
		p.carts[cartId][topic] = ch
		logger.Info("subscribed", log.Fields{"topic": topic, "cart_id": cartId})
		return ch
	}
	return ch
}
