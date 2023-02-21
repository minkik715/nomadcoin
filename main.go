package main

import (
	"github.com/rlaalsrl715/nomadcoin/explorer"
	"github.com/rlaalsrl715/nomadcoin/rest"
)

func main() {
	go rest.Start(8080)
	explorer.Start(8000)
}
