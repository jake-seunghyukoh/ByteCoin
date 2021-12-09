package cli

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/ohshyuk5/ByteCoin/explorer"
	"github.com/ohshyuk5/ByteCoin/rest"
)

func usage() {
	fmt.Println("Welcome to Byte Coin")
	fmt.Println()
	fmt.Println("Please use the following flags:")
	fmt.Println()
	fmt.Println("-port:		Set port of the server")
	fmt.Println("-mode:		Choose between 'html' and 'rest'")
	fmt.Println()
	runtime.Goexit()
}

func Start() {
	if len(os.Args) == 1 {
		usage()
	}

	port := flag.Int("port", 4000, "Set port of the server")
	mode := flag.String("mode", "rest", "Choose between 'html' and 'rest'")

	flag.Parse()

	switch *mode {
	case "rest":
		// start rest api
		rest.Start(*port)
	case "html":
		// start html explorer
		explorer.Start(*port)
	default:
		usage()
	}

	fmt.Println(*port, *mode)
}
