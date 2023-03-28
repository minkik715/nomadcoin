package wallet

import (
	"crypto/ecdsa"
	"os"
)

type wallet struct {
	privateKey *ecdsa.PrivateKey
}

var w *wallet

func hasWalletFile() bool {
	_, err := os.Stat("nico-coin.wallet")
	return !os.IsNotExist(err)
}

func Wallet() *wallet {
	if w == nil {
		if hasWalletFile() {

		} else {

		}
	}
}
