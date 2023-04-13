package main

import (
	"github.com/rlaalsrl715/nomadcoin/cli"
	"github.com/rlaalsrl715/nomadcoin/db"
)

func main() {
	defer db.DB().Close()

	cli.Start()
}
