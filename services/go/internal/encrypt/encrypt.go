package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"github.com/google/tink/go/daead/subtle"
	"hash"
	"io"
)

// AssociatedData
// TODO load this from k8s
var AssociatedData = []byte("sa_temporal-shop")

// SHA256 hashes using sha256 algorithm
func SHA256(text string) string {
	algorithm := sha256.New()
	return stringHasher(algorithm, text)
}
func stringHasher(algorithm hash.Hash, text string) string {
	algorithm.Write([]byte(text))
	return hex.EncodeToString(algorithm.Sum(nil))
}
func MDHash(input string) string {
	byteInput := []byte(input)

	/* we arent concerned about collision resistance in this */
	/* #nosec */
	md5Hash := md5.Sum(byteInput)
	return hex.EncodeToString(md5Hash[:]) // by referring to it as a string
}
func Encrypt(key string, value []byte) ([]byte, error) {

	aesBlock, err := aes.NewCipher([]byte(MDHash(key)))
	if err != nil {
		return nil, err
	}
	// gcm or Galois/Counter Mode, is a mode of operation
	// for symmetric key cryptographic block ciphers
	// - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	gcm, err := cipher.NewGCM(aesBlock)
	// if any error generating new GCM
	// handle them
	if err != nil {
		return nil, err
	}
	// creates a new byte array the size of the nonce
	// which must be passed to Seal
	nonce := make([]byte, gcm.NonceSize())
	// populates our nonce with a cryptographically secure
	// random sequence
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	out := gcm.Seal(nonce, nonce, value, nil)
	return out, nil
}

func Decrypt(key string, value []byte) ([]byte, error) {
	hashedPhrase := MDHash(key)
	aesBlock, err := aes.NewCipher([]byte(hashedPhrase))
	if err != nil {
		return nil, err
	}
	gcmInstance, err := cipher.NewGCM(aesBlock)
	if err != nil {
		return nil, err
	}
	nonceSize := gcmInstance.NonceSize()
	nonce, cipheredText := value[:nonceSize], value[nonceSize:]
	out, err := gcmInstance.Open(nil, nonce, cipheredText, nil)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func EncryptDeterministically(key []byte, value []byte, associatedData []byte) (string, error) {
	hashedKey := SHA256(string(key))
	sut, err := subtle.NewAESSIV([]byte(hashedKey))
	if err != nil {
		return "", err
	}
	actual, err := sut.EncryptDeterministically(value, associatedData)
	if err != nil {
		return "", err
	}
	encoded := hex.EncodeToString(actual)
	return encoded, nil
}
func DecryptDeterministically(key []byte, value []byte, associatedData []byte) ([]byte, error) {
	hashedKey := SHA256(string(key))
	sut, err := subtle.NewAESSIV([]byte(hashedKey))
	if err != nil {
		return nil, err
	}
	actual, err := sut.DecryptDeterministically(value, associatedData)
	return actual, err

}
