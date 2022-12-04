package shopping

import (
	"github.com/temporalio/temporal-shop/services/go/internal/encrypt"
)

var MDHash = encrypt.MDHash

func GenerateShopperHash(key, email string) (string, error) {
	return email, nil
}
func ExtractShopperEmail(key, value string) (string, error) {
	return value, nil
}
