package acmewrapper

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"os"
)

// SavePrivateKey is used to write the given key to file
// This code copied verbatim from caddy:
// https://github.com/mholt/caddy/blob/master/caddy/https/crypto.go
func SavePrivateKey(filename string, key crypto.PrivateKey) error {
	var pemType string
	var keyBytes []byte
	switch key := key.(type) {
	case *ecdsa.PrivateKey:
		var err error
		pemType = "EC"
		keyBytes, err = x509.MarshalECPrivateKey(key)
		if err != nil {
			return err
		}
	case *rsa.PrivateKey:
		pemType = "RSA"
		keyBytes = x509.MarshalPKCS1PrivateKey(key)
	}

	pemKey := pem.Block{Type: pemType + " PRIVATE KEY", Bytes: keyBytes}
	keyOut, err := os.Create(filename)
	if err != nil {
		return err
	}
	keyOut.Chmod(0600)
	defer keyOut.Close()
	return pem.Encode(keyOut, &pemKey)
}

// LoadPrivateKey reads a key from file
// This code copied verbatim from caddy:
// https://github.com/mholt/caddy/blob/master/caddy/https/crypto.go
func LoadPrivateKey(filename string) (crypto.PrivateKey, error) {
	keyBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	keyBlock, _ := pem.Decode(keyBytes)

	switch keyBlock.Type {
	case "RSA PRIVATE KEY":
		return x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
	case "EC PRIVATE KEY":
		return x509.ParseECPrivateKey(keyBlock.Bytes)
	}

	return nil, errors.New("unknown private key type")
}
