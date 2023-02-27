package cli

import (
	"flag"
	"fmt"
	"github.com/rlaalsrl715/nomadcoin/explorer"
	"github.com/rlaalsrl715/nomadcoin/rest"
	"os"
)

func Start() {
	if len(os.Args) < 2 {
		usage()
	}
	var port = flag.Int("port", 4000, "Set port of the server")
	var mode = flag.String("mode", "rest", "Choose between 'html' and 'rest'")

	flag.Parse()

	switch *mode {
	case "rest":
		rest.Start(*port)
	case "html":
		explorer.Start(*port)
	default:
		usage()
	}
}

func usage() {
	fmt.Printf("Welcom to nomad coin\n\n")
	fmt.Printf("Please use the follwing flags\n\n")
	fmt.Printf("-port=4000:  Set the PORT of the server\n")
	fmt.Printf("-mode:  Choose between 'html' and 'rest'(recommended rest)\n")
	os.Exit(0)
}
