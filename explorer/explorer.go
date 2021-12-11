package explorer

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/ohshyuk5/ByteCoin/utils"

	"github.com/ohshyuk5/ByteCoin/blockchain"
)

const (
	templateDir string = "explorer/templates/"
)

var port string = ":4000"

var templates *template.Template

type homeData struct {
	PageTitle string
	Blocks    []*blockchain.Block
}

func home(rw http.ResponseWriter, _ *http.Request) {
	data := homeData{PageTitle: "Home", Blocks: blockchain.Blocks(blockchain.BlockChain())}
	err := templates.ExecuteTemplate(rw, "home", data)
	utils.HandleErr(err)
}

func add(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := templates.ExecuteTemplate(rw, "add", nil)
		utils.HandleErr(err)
	case "POST":
		err := r.ParseForm()
		utils.HandleErr(err)
		blockchain.BlockChain().AddBlock()
		http.Redirect(rw, r, "/", http.StatusPermanentRedirect)
	}
}

func Start(aPort int) {
	port = fmt.Sprintf(":%d", aPort)

	templates = template.Must(template.ParseGlob(templateDir + "pages/*.gohtml"))
	templates = template.Must(templates.ParseGlob(templateDir + "partials/*.gohtml"))

	handler := http.NewServeMux()
	handler.HandleFunc("/", home)
	handler.HandleFunc("/add", add)

	fmt.Printf("Listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, handler))
}
