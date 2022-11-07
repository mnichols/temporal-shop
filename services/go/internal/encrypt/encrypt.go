package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"io"
)

func mdHash(input string) string {
	byteInput := []byte(input)
	md5Hash := md5.Sum(byteInput)
	return hex.EncodeToString(md5Hash[:]) // by referring to it as a string
}
func Encrypt(key string, value []byte) ([]byte, error) {

	aesBlock, err := aes.NewCipher([]byte(mdHash(key)))
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
	hashedPhrase := mdHash(key)
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
