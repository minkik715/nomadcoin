package main

import (
	"fmt"
	"github.com/rlaalsrl715/nomadcoin/blockchian"
	"html/template"
	"log"
	"net/http"
)

const port string = ":4000"

type homeData struct {
	PageTitle string
	Blocks    []*blockchian.Block
}

func home(rw http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("template/home.gohtml"))
	data := homeData{"Welcome", blockchian.GetBlockchain().AllBLocks()}
	tmpl.Execute(rw, data)
}

func main() {
	http.HandleFunc("/", home)
	fmt.Printf("Listening on http://localhost%s\n", port)
	blockchian.GetBlockchain().AppendBlock("NICO!!!!!!")

	log.Fatal(http.ListenAndServe(port, nil))
}
