package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"github.com/rlaalsrl715/nomadcoin/utils"
	"math/big"
	"os"
)

const (
	walletName string = "nomadcoin.wallet"
)

type wallet struct {
	privateKey *ecdsa.PrivateKey
	Address    string
}

var w *wallet

func hasWalletFile() bool {
	_, err := os.Stat(walletName)
	return !os.IsNotExist(err)
}

func createPrivateKey() *ecdsa.PrivateKey {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	utils.HandleErr(err)
	return privateKey
}

func persistKey(key *ecdsa.PrivateKey) {
	bytes, err := x509.MarshalECPrivateKey(key)
	utils.HandleErr(err)
	err = os.WriteFile(walletName, bytes, 0644)
	utils.HandleErr(err)
}

func restoreKey() (key *ecdsa.PrivateKey) {
	bytes, err := os.ReadFile(walletName)
	utils.HandleErr(err)
	key, err = x509.ParseECPrivateKey(bytes)
	utils.HandleErr(err)
	return
}

func aFromKey(key *ecdsa.PrivateKey) string {
	return encodingBytes(key.X.Bytes(), key.Y.Bytes())
}

func encodingBytes(a, b []byte) string {
	bytes := append(a, b...)
	payload := fmt.Sprintf("%x", bytes)
	return payload
}

func Sign(payload string, w *wallet) string {
	bytes, err := hex.DecodeString(payload)
	utils.HandleErr(err)
	r, s, err := ecdsa.Sign(rand.Reader, w.privateKey, bytes)
	utils.HandleErr(err)
	return encodingBytes(r.Bytes(), s.Bytes())
}

func Verify(signature string, hashMessage string, address string) bool {
	r, s, err := restoreBigInts(signature)
	utils.HandleErr(err)
	x, y, err := restoreBigInts(address)
	utils.HandleErr(err)
	hashMessageB, err := hex.DecodeString(hashMessage)
	utils.HandleErr(err)
	publicKey := ecdsa.PublicKey{Curve: elliptic.P256(), X: x, Y: y}
	return ecdsa.Verify(&publicKey, hashMessageB, r, s)

}

func restoreBigInts(payload string) (*big.Int, *big.Int, error) {
	signatureB, err := hex.DecodeString(payload)
	if err != nil {
		return nil, nil, err
	}
	bigA, bigB := big.Int{}, big.Int{}
	bigA.SetBytes(signatureB[:len(signatureB)/2])
	bigB.SetBytes(signatureB[len(signatureB)/2:])
	return &bigA, &bigB, nil
}

func Wallet() *wallet {
	if w == nil {
		w = &wallet{}
		if hasWalletFile() {
			key := restoreKey()
			w.privateKey = key
		} else {
			key := createPrivateKey()
			persistKey(key)
			w.privateKey = key
		}
		w.Address = aFromKey(w.privateKey)
	}
	return w
}
