package shopping

var TypeHandlers *Handlers

type Handlers struct{}

func NewHandlers() *Handlers {
	return &Handlers{}
}
