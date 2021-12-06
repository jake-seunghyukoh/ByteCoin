package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/ohshyuk5/ByteCoin/blockchain"
)

const (
	port        string = ":4000"
	templateDir string = "templates/"
)

var templates *template.Template

type homeData struct {
	PageTitle string
	Blocks    []*blockchain.Block
}

func home(rw http.ResponseWriter, r *http.Request) {
	data := homeData{PageTitle: "Home", Blocks: blockchain.GetBlockChain().AllBlocks()}
	templates.ExecuteTemplate(rw, "home", data)
}

func main() {
	templates = template.Must(template.ParseGlob(templateDir + "pages/*.gohtml"))
	templates = template.Must(templates.ParseGlob(templateDir + "partials/*.gohtml"))

	http.HandleFunc("/home", home)
	fmt.Printf("Listening on http://localhost%s/home\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
