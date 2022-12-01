package encrypt

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncryptAndDecrypt(t *testing.T) {
	key := "foo123"
	email := "mike@example.com"

	actual, err := Encrypt(key, []byte(email))
	if err != nil {
		t.Fatalf("generating hash should not have error %s", err)
	}
	if len(actual) == 0 {
		t.Fatal("actual hash must not be empty")
	}
	decrypted, err := Decrypt(key, actual)
	if err != nil {
		t.Fatalf("descrypting should work %s", err)
	}
	if string(decrypted) != email {
		t.Errorf("expected '%s' but got '%s'", email, decrypted)
	}
}
func TestEncryptReuse(t *testing.T) {
	A := assert.New(t)
	key := "foo123"
	email := "mike@example.com"

	expect, err := Encrypt(key, []byte(email))
	A.NoError(err)

	actual, err := Encrypt(key, []byte(email))
	A.NoError(err)
	A.Equal(expect, actual)
}
