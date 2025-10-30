package cursor

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
)

const (
	cursorContent = iota
	cursorSignature
	cursorLen
)

// Decrypt decrypts the cursor, ensures its integrity by verifying its HMAC signature.
func Decrypt[T Pointer](content, secret []byte) (*Cursor[T], error) {
	raw := bytes.Split(content, sep)
	if len(raw) != cursorLen {
		return nil, fmt.Errorf("parsing: invalid cursor format")
	}
	mac := hmac.New(sha256.New, secret)
	_, err := mac.Write(raw[cursorContent])
	if err != nil {
		return nil, fmt.Errorf("hashing: %w", err)
	}
	if !hmac.Equal(mac.Sum(nil), raw[cursorSignature]) {
		return nil, fmt.Errorf("checking: %w", err)
	}
	var c Cursor[T]
	err = c.Decode(raw[cursorContent])
	if err != nil {
		return nil, fmt.Errorf("unmarshaling: %w", err)
	}
	return &c, nil
}

// Encrypt encrypts the cursor and then ensures it integrity by signing the content.
// It concatenates the base64-encoded JSON representation of the cursor with a dot and its sha256 hmac signature.
func Encrypt[T Pointer](c *Cursor[T], secret []byte) ([]byte, error) {
	src, err := c.Encode()
	if err != nil {
		return nil, fmt.Errorf("marshalling: %w", err)
	}
	if len(src) == 0 {
		return nil, nil
	}
	mac := hmac.New(sha256.New, secret)
	_, err = mac.Write(src)
	if err != nil {
		return nil, fmt.Errorf("hashing: %w", err)
	}
	sum := mac.Sum(nil)
	dst := make([]byte, b64.EncodedLen(len(sum)))
	b64.Encode(dst, sum)

	return bytes.Join([][]byte{src, dst}, sep), nil
}
