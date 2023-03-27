package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"github.com/rlaalsrl715/nomadcoin/utils"
)

func Start() {
	message := "fuck github"
	hashMessage := utils.Hash(message)

	hashByteMessage, err := hex.DecodeString(hashMessage)
	utils.HandleErr(err)

	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	utils.HandleErr(err)

	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hashByteMessage)
	utils.HandleErr(err)
}
