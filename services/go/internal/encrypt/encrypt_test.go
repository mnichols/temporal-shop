package encrypt

import (
	"encoding/hex"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"github.com/temporalio/temporal-shop/api/temporal_shop/test/v1"
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
	A.NotEqual(expect, actual)
}

// https://pkg.go.dev/github.com/google/tink/go/daead/subtle
func TestDeterministicEncryption(t *testing.T) {
	A := assert.New(t)
	key := []byte("foo")
	sessID := &test.TestSessionID{Email: "myemail@example.org"}
	out, err := proto.Marshal(sessID)
	A.NoError(err)
	dataYouWantToEncrypt := out
	authenticatedDataThatIsASourceYouTrust := "myhostname.com"
	encrypted, err := EncryptDeterministically(key, out, []byte(authenticatedDataThatIsASourceYouTrust))
	A.NoError(err)
	fmt.Println("sessionID", encrypted)

	decoded, err := hex.DecodeString(encrypted)
	A.NoError(err)
	decrypted, err := DecryptDeterministically(key, decoded, []byte(authenticatedDataThatIsASourceYouTrust))
	A.NoError(err)
	A.Equal(dataYouWantToEncrypt, decrypted)
	encryptedAgain, err := EncryptDeterministically(key, dataYouWantToEncrypt, []byte(authenticatedDataThatIsASourceYouTrust))
	A.NoError(err)
	fmt.Println("sessionIDagain", encrypted)

	A.Equal(encrypted, encryptedAgain)
	different, err := EncryptDeterministically(key, dataYouWantToEncrypt, []byte("somethingelse"))
	A.NoError(err)
	A.NotEqual(encrypted, different)

}
