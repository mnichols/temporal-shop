package shopping

import (
	"github.com/temporalio/temporal-shop/services/go/internal/encrypt"
)

var MDHash = encrypt.MDHash

func GenerateShopperHash(key, email string) (string, error) {
	return email, nil
	//value, err := encrypt.Encrypt(key, []byte(email))
	//if err != nil {
	//	return "", err
	//}
	//return hex.EncodeToString(value), nil
}
func ExtractShopperEmail(key, value string) (string, error) {
	return value, nil
	//decoded, err := hex.DecodeString(value)
	//if err != nil {
	//	return "", err
	//}
	//result, err := encrypt.Decrypt(key, decoded)
	//if err != nil {
	//	return "", err
	//}
	//return string(result), nil
}
