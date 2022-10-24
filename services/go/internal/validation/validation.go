package validation

import (
	"context"

	"sync"

	"github.com/go-playground/validator/v10"
)

var once sync.Once

var singletonValidate *validator.Validate

func MustGetValidator(ctx context.Context) *validator.Validate {
	once.Do(func() {
		singletonValidate = validator.New()
	})
	return singletonValidate
}
