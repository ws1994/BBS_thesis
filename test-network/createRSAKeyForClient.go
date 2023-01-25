package main

import (
	// "crypto/aes"
	// "crypto/cipher"
	"crypto/rand"
	// "crypto/sha256"
	// "encoding/hex"
	// "io"
	// "log"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
)

func InitRSAKeyForCC() error {
	// peerMSP, _ := shim.GetMSPID()
	rsaPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	// ***************************
	// write private key to file
	// ***************************
	objPkcs8, _ := x509.MarshalPKCS8PrivateKey(rsaPrivateKey)
	privateKey := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: objPkcs8,
	}
	file, err := os.Create("./clientRSAKeys/ccPrivate.pem")
	if err != nil {
		return err
	}
	err = pem.Encode(file, privateKey)

	// ***************************
	// write public key to file
	// ***************************
	objPkix, err := x509.MarshalPKIXPublicKey(&rsaPrivateKey.PublicKey)
	if err != nil {
		return err
	}
	publicKey := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: objPkix,
	}
	file2, err := os.Create("./clientRSAKeys/ccPublic.pem")
	if err != nil {
		return err
	}
	err = pem.Encode(file2, publicKey)
	return nil
}

func main() {
	InitRSAKeyForCC()
}
