package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/ohshyuk5/ByteCoin/blockchain"
)

const port string = ":4000"

type homeData struct {
	PageTitle string
	Blocks    []*blockchain.Block
}

func home(rw http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/home.html"))
	tmpl.Execute(rw, homeData{PageTitle: "Home", Blocks: blockchain.GetBlockChain().AllBlocks()})
}

func main() {
	http.HandleFunc("/home", home)
	fmt.Printf("Listening on http://localhost%s/home\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
