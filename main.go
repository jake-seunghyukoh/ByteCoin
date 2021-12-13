package main

import (
	"github.com/ohshyuk5/ByteCoin/cli"
	"github.com/ohshyuk5/ByteCoin/db"
)

func main() {
	defer db.Close()
	cli.Start()
}
