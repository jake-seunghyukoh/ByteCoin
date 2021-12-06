package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ohshyuk5/ByteCoin/rest"
)

func main() {
	http.HandleFunc("/", rest.Documentation)
	http.HandleFunc("/blocks", rest.Blocks)
	fmt.Printf("Listening on %s%s\n", rest.BaseURL, rest.Port)
	log.Fatal(http.ListenAndServe(rest.Port, nil))
}
