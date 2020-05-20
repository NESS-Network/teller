package main

import (
	"github.com/MDLlife/MDL/src/api"
	"os"
	"log"
)

func main() {

	if len(os.Args) < 2 {
		println("Usage: addr <wallet_filanem>")
		return
	}

	client := api.NewClient("http://127.0.0.1:6420")

	walletResp, err := client.Wallet(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}

	for _, wallet := range walletResp.Entries {
		println(wallet.Address)
	}
}
