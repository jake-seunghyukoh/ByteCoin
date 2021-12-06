package main

import (
	"github.com/ohshyuk5/ByteCoin/explorer"
	"github.com/ohshyuk5/ByteCoin/rest"
)

func main() {
	go explorer.Start(3000)
	rest.Start(4000)
}
