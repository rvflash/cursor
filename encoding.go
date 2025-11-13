// Copyright (c) 2025 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package cursor

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"fmt"
)

const (
	cursorContent = iota
	cursorSignature
	_cursorLen
)

// Decrypt decrypts the cursor, ensures its integrity by verifying its HMAC signature.
func Decrypt[T Pointer](content, secret []byte) (*Cursor[T], error) {
	raw := bytes.Split(content, sep)
	if len(raw) != _cursorLen {
		return nil, fmt.Errorf("parsing: invalid cursor format")
	}
	src, err := b64Decode(raw[cursorSignature])
	if err != nil {
		return nil, fmt.Errorf("hash decoding: %w", err)
	}
	sig, err := sign(raw[cursorContent], secret)
	if err != nil {
		return nil, fmt.Errorf("signature checking: %w", err)
	}
	if !hmac.Equal(src, sig) {
		return nil, errors.New("signature mismatch")
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
	sig, err := sign(src, secret)
	if err != nil {
		return nil, fmt.Errorf("signing: %w", err)
	}
	return bytes.Join([][]byte{src, b64Encode(sig)}, sep), nil
}

func sign(content, secret []byte) ([]byte, error) {
	mac := hmac.New(sha256.New, secret)
	_, err := mac.Write(content)
	if err != nil {
		return nil, err
	}
	return mac.Sum(nil), nil
}
