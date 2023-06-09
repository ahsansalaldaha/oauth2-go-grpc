package services

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/golang-jwt/jwt"
)

var privateKey *rsa.PrivateKey
var publicKey *rsa.PublicKey

func initKeys(path string) error {
	// Load the public key from a PEM-encoded file
	publicKeyData, err := ioutil.ReadFile(path + "/public_key.pem")
	if err != nil {
		return fmt.Errorf("Failed to load public key: %v", err)
	}

	publicKey, err = jwt.ParseRSAPublicKeyFromPEM(publicKeyData)
	if err != nil {
		return fmt.Errorf("Failed to parse public key: %v", err)
	}
	// Load the private key from a PEM-encoded file
	privateKeyData, err := ioutil.ReadFile(path + "/private_key.pem")
	if err != nil {
		return fmt.Errorf("Failed to load private key: %v", err)
	}

	privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateKeyData)
	if err != nil {
		return fmt.Errorf("Failed to parse private key: %v", err)
	}
	return nil
}

// GetPrivateKey - Returns the private key generated
func GetPrivateKey() *rsa.PrivateKey {
	return privateKey
}

// GetPublicKey - Returns the private key generated
func GetPublicKey() *rsa.PublicKey {
	return publicKey
}

// GenerateRSAKeyPair -  Generates RSA key pair of certain bit size
func GenerateRSAKeyPair(bitSize int, privateKeyFile, publicKeyFile string) error {
	// Generate a private key
	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return fmt.Errorf("failed to generate private key: %v", err)
	}

	// Save the private key in PEM format
	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	privateKeyFileWriter, err := os.Create(privateKeyFile)
	if err != nil {
		return fmt.Errorf("failed to create private key file: %v", err)
	}
	defer privateKeyFileWriter.Close()

	err = pem.Encode(privateKeyFileWriter, privateKeyPEM)
	if err != nil {
		return fmt.Errorf("failed to save private key: %v", err)
	}

	// Save the public key in PEM format
	publicKey := &privateKey.PublicKey
	publicKeyDER, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return fmt.Errorf("failed to marshal public key: %v", err)
	}

	publicKeyPEM := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyDER,
	}

	publicKeyFileWriter, err := os.Create(publicKeyFile)
	if err != nil {
		return fmt.Errorf("failed to create public key file: %v", err)
	}
	defer publicKeyFileWriter.Close()

	err = pem.Encode(publicKeyFileWriter, publicKeyPEM)
	if err != nil {
		return fmt.Errorf("failed to save public key: %v", err)
	}

	return nil
}

// GenerateRSAKeyPairIfNotExists - Generate RSA key Pair if they are not already generated
func GenerateRSAKeyPairIfNotExists(path string) error {
	if _, err := os.Stat(path + "/private_key.pem"); err != nil {
		err := GenerateRSAKeyPair(2048, path+"/private_key.pem", path+"/public_key.pem")
		if err != nil {
			panic(fmt.Sprintf("Failed to generate RSA key pair: %v", err))

		}
	}
	return initKeys(path)
}
