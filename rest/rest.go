package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ohshyuk5/ByteCoin/blockchain"
	"github.com/ohshyuk5/ByteCoin/utils"
)

const baseURL string = "http://localhost"

var port string = ":4000"

type url string

func (u url) MarshalText() ([]byte, error) {
	url := fmt.Sprintf("%s%s%s", baseURL, port, u)
	return []byte(url), nil
}

type urlDescription struct {
	URL         url    `json:"url"`
	Method      string `json:"method"`
	Description string `json:"description"`
	Payload     string `json:"payload,omitempty"`
}

type addBlockBody struct {
	Message string
}

func documentation(rw http.ResponseWriter, r *http.Request) {
	data := []urlDescription{
		{
			URL:         url("/"),
			Method:      "GET",
			Description: "See documentation",
		},
		{
			URL:         url("/blocks"),
			Method:      "POST",
			Description: "Add a Block",
			Payload:     "data: string",
		},
		{
			URL:         url("/blocks"),
			Method:      "GET",
			Description: "See all Blocks",
		},
		{
			URL:         url("/blocks/{id}"),
			Method:      "GET",
			Description: "See a Block",
		},
	}
	rw.Header().Add("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(data)
}

func blocks(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		rw.Header().Add("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(blockchain.GetBlockChain().AllBlocks())
	case "POST":
		var body addBlockBody
		err := json.NewDecoder(r.Body).Decode(&body)
		utils.HandleErr(err)

		blockchain.GetBlockChain().AddBlock(body.Message)
		rw.WriteHeader(http.StatusCreated)
	}
}

func Start(aPort int) {
	port = fmt.Sprintf(":%d", aPort)

	handler := http.NewServeMux()
	handler.HandleFunc("/", documentation)
	handler.HandleFunc("/blocks", blocks)

	fmt.Printf("Listening on %s%s\n", baseURL, port)
	log.Fatal(http.ListenAndServe(port, handler))
}
