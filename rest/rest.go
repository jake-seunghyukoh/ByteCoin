package rest

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ohshyuk5/ByteCoin/blockchain"
	"github.com/ohshyuk5/ByteCoin/utils"
	"log"
	"net/http"
)

const baseURL string = "http://localhost"

var port = ":4000"

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

type balanceResponse struct {
	Address string `json:"address"`
	Balance int    `json:"balance"`
}

type errorResponse struct {
	ErrorMessage string `json:"errorMessage"`
}

type addTxPayload struct {
	To     string `json:"to"`
	Amount int    `json:"amount"`
}

func documentation(rw http.ResponseWriter, _ *http.Request) {
	data := []urlDescription{
		{
			URL:         url("/"),
			Method:      "GET",
			Description: "See documentation",
		},
		{
			URL:         url("/status"),
			Method:      "GET",
			Description: "See the status of the blockchain",
		},
		{
			URL:         url("/blocks"),
			Method:      "POST",
			Description: "Add a Block",
		},
		{
			URL:         url("/blocks"),
			Method:      "GET",
			Description: "See all Blocks",
		},
		{
			URL:         url("/blocks/{hash}"),
			Method:      "GET",
			Description: "See a Block",
		},
		{
			URL:         url("/balance/{address}"),
			Method:      "GET",
			Description: "Get TxOuts for an Address",
		},
	}
	utils.HandleErr(json.NewEncoder(rw).Encode(data))
}

func blocks(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		utils.HandleErr(json.NewEncoder(rw).Encode(blockchain.BlockChain().Blocks()))
	case "POST":
		blockchain.BlockChain().AddBlock()
		rw.WriteHeader(http.StatusCreated)
	}
}

func block(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]
	block, err := blockchain.FindBlock(hash)
	encoder := json.NewEncoder(rw)

	if err == blockchain.ErrNotFound {
		rw.WriteHeader(http.StatusNotFound)
		utils.HandleErr(encoder.Encode(errorResponse{ErrorMessage: fmt.Sprint(err)}))
	} else {
		utils.HandleErr(encoder.Encode(block))
	}
}

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(rw, r)
	})
}

func status(rw http.ResponseWriter, _ *http.Request) {
	err := json.NewEncoder(rw).Encode(blockchain.BlockChain())
	utils.HandleErr(err)
}

func balance(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]
	total := r.URL.Query().Get("total")
	switch total {
	case "true":
		balance := blockchain.BlockChain().BalanceByAddress(address)
		utils.HandleErr(
			json.NewEncoder(rw).Encode(balanceResponse{
				Address: address,
				Balance: balance,
			}),
		)
	default:
		utils.HandleErr(
			json.NewEncoder(rw).Encode(
				blockchain.BlockChain().TxOutsByAddress(address),
			),
		)
	}
}

func mempool(rw http.ResponseWriter, _ *http.Request) {
	utils.HandleErr(
		json.NewEncoder(rw).Encode(blockchain.Mempool.Txs),
	)
}

func transactions(rw http.ResponseWriter, r *http.Request) {
	var payload addTxPayload
	utils.HandleErr(json.NewDecoder(r.Body).Decode(&payload))

	err := blockchain.Mempool.AddTx(payload.To, payload.Amount)
	if err != nil {
		utils.HandleErr(json.NewEncoder(rw).Encode(errorResponse{"not enough funds"}))
	}
	rw.WriteHeader(http.StatusCreated)
}

func Start(aPort int) {
	port = fmt.Sprintf(":%d", aPort)

	router := mux.NewRouter()

	router.Use(jsonContentTypeMiddleware)

	router.HandleFunc("/", documentation).Methods("GET")
	router.HandleFunc("/status", status).Methods("GET")
	router.HandleFunc("/blocks", blocks).Methods("GET", "POST")
	router.HandleFunc("/blocks/{hash:[a-f0-9]+}", block).Methods("GET")
	router.HandleFunc("/balance/{address}", balance).Methods("GET")
	router.HandleFunc("/mempool", mempool).Methods("GET")
	router.HandleFunc("/transactions", transactions).Methods("POST")

	fmt.Printf("Listening on %s%s\n", baseURL, port)
	log.Fatal(http.ListenAndServe(port, router))
}
