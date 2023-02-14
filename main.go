package main

import (
	"fmt"
	"github.com/rlaalsrl715/nomadcoin/blockchian"
)

func main() {
	chain := blockchian.GetBlockchain()
	chain.AppendBlock("nico")
	chain.AppendBlock("nico2")
	chain.AppendBlock("nico23")
	chain.AppendBlock("nico234")
	for _, b := range chain.AllBLocks() {
		fmt.Println(b)
	}
}
