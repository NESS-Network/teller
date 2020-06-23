package main

import (
	"os"
	"log"
	"fmt"
	"bufio"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/common"
	"context"
	"math/big"
	"math"
	"strings"
)

const MAX_COUNT  = 100

func main() {
	client, err := ethclient.Dial("http://127.0.0.1:8545")
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) < 2 {
		println("Usage: balance <addr_list_file>")
		return
	}

	var addrs []string
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	ethTotalValue := new(big.Float)

	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() {
		addr := scanner.Text()
		addrs = append(addrs, addr)

		account := common.HexToAddress(addr)
		balance, err := client.BalanceAt(context.Background(), account, nil)
		if err != nil {
			log.Fatal(err)
		}
		fbalance := new(big.Float)
		fbalance.SetString(balance.String())
		ethValue := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))
		ethTotalValue = ethTotalValue.Add(ethTotalValue, ethValue)
		//println(ethValue.Text('f',7))
		if strings.Compare(ethValue.Text('f',8), "0.00000000") != 0 {
			println(addr, fmt.Sprint(ethValue))
		} else {
			pendingBalance, err := client.PendingBalanceAt(context.Background(), account)
			if err != nil {
				log.Fatal(err)
			}
			fbalance = new(big.Float)
			fbalance.SetString(pendingBalance.String())
			pendingBalanceEthValue := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))
			if strings.Compare(pendingBalanceEthValue.Text('f', 8), "0.00000000") != 0 {
				println(addr, fmt.Sprint(pendingBalanceEthValue), "(pending)")
			}
		}
		count++
		if count > MAX_COUNT {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	println("Addresses total balance:", fmt.Sprint(ethTotalValue))




}
