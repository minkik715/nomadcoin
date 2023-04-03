package main

import (
	"github.com/rlaalsrl715/nomadcoin/db"
	"github.com/rlaalsrl715/nomadcoin/rest"
)

func main() {
	defer db.DB().Close()

	rest.Start(3000)
}
