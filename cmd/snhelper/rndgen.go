package main

import (
	//	"crypto/ecdsa"
	//	"crypto/elliptic"
	//	"crypto/rand"
	"encoding/hex"
	//	"fmt"

	//	"github.com/ovcharovvladimir/essentiaHybrid/contracts/ens/contract"
	"github.com/ovcharovvladimir/essentiaHybrid/crypto"
)

func Key() (string, error) {

	key, err := crypto.GenerateKey()
	if err != nil {
		//utils.Fatalf("Failed to generate private key: %s", err)
		return "", err
	}
	k := hex.EncodeToString(crypto.FromECDSA(key))

	return k, nil
}
