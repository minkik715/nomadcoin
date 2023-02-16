package explorer

import (
	"fmt"
	"github.com/rlaalsrl715/nomadcoin/blockchian"
	"html/template"
	"log"
	"net/http"
)

const (
	port        string = ":4000"
	templateDir string = "explorer/templates/"
)

var templates *template.Template

type homeData struct {
	PageTitle string
	Blocks    []*blockchian.Block
}

func home(rw http.ResponseWriter, r *http.Request) {
	data := homeData{"Home", blockchian.GetBlockchain().AllBLocks()}
	templates.ExecuteTemplate(rw, "home", data)
}

func add(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		templates.ExecuteTemplate(rw, "add", nil)
	case "POST":
		r.ParseForm()
		data := r.Form.Get("blockData")
		blockchian.GetBlockchain().AppendBlock(data)
		http.Redirect(rw, r, "/", http.StatusPermanentRedirect)

	}
	templates.ExecuteTemplate(rw, "add", nil)
}

func Start() {
	templates = template.Must(template.ParseGlob(templateDir + "pages/*.gohtml"))
	templates = template.Must(templates.ParseGlob(templateDir + "partials/*.gohtml"))
	http.HandleFunc("/", home)
	http.HandleFunc("/add", add)
	fmt.Printf("Listening on http://localhost%s\n", port)
	blockchian.GetBlockchain().AppendBlock("NICO!!!!!!")

	log.Fatal(http.ListenAndServe(port, nil))
}
