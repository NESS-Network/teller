package main

import (
	"os"
	"log"
	"fmt"
	"strings"
	"bufio"
	"github.com/MDLlife/MDL/src/api"
)

func main() {

	if len(os.Args) < 2 {
		println("Usage: addr <wallet_filename> / <addr_list_file>")
		return
	}

	client := api.NewClient("http://127.0.0.1:6420")

	if strings.Contains(os.Args[1], ".wlt") {
		wb, err := client.WalletBalance(os.Args[1])
		if err != nil {
			log.Fatalln(err)
		}
		println("Wallet Balance:", fmt.Sprint(wb.Confirmed.Coins))
	} else {
		var addrs []string
		file, err := os.Open(os.Args[1])
		if err != nil {
			log.Fatalln(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			addrs = append(addrs, scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			log.Fatalln(err)
		}

		br, err := client.Balance(addrs)
		if err != nil {
			log.Fatalln(err)
		}
		println("Addresses total balance:", fmt.Sprint(br.Confirmed.Coins))
	}



}
