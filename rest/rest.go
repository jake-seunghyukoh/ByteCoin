package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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

type errorResponse struct {
	ErrorMessage string `json:"errorMessage"`
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
			URL:         url("/blocks/{height}"),
			Method:      "GET",
			Description: "See a Block",
		},
	}
	json.NewEncoder(rw).Encode(data)
}

func blocks(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		json.NewEncoder(rw).Encode(blockchain.GetBlockChain().AllBlocks())
	case "POST":
		var body addBlockBody
		err := json.NewDecoder(r.Body).Decode(&body)
		utils.HandleErr(err)

		blockchain.GetBlockChain().AddBlock(body.Message)
		rw.WriteHeader(http.StatusCreated)
	}
}

func block(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	heightStr := vars["height"]
	heightInt, err := strconv.Atoi(heightStr)
	utils.HandleErr(err)

	block, err := blockchain.GetBlockChain().GetBlock(heightInt)
	encoder := json.NewEncoder(rw)

	if err == blockchain.ErrNotFound {
		rw.WriteHeader(http.StatusNotFound)
		encoder.Encode(errorResponse{ErrorMessage: fmt.Sprint(err)})
	} else {
		encoder.Encode(block)
	}
}

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(rw, r)
	})
}

func Start(aPort int) {
	port = fmt.Sprintf(":%d", aPort)

	router := mux.NewRouter()

	router.Use(jsonContentTypeMiddleware)

	router.HandleFunc("/", documentation).Methods("GET")
	router.HandleFunc("/blocks", blocks).Methods("GET", "POST")
	router.HandleFunc("/blocks/{height:[0-9]+}", block).Methods("GET")

	fmt.Printf("Listening on %s%s\n", baseURL, port)
	log.Fatal(http.ListenAndServe(port, router))
}
