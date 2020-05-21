package main

import (
	"crypto/ecdsa"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"os"
)

func main() {
	prFile, err := os.Create("eth_list_pr_key.txt")
	if err != nil {
		log.Fatalln(err)
	}
	defer prFile.Close()
	puFile, err := os.Create("eth_list_pu_key.txt")
	if err != nil {
		log.Fatalln(err)
	}
	defer puFile.Close()

	for i := range [10000]int{} {
		pr, pu := generatePrPuKeyPair()
		fmt.Println(i, pr, pu)
		prFile.WriteString(pr+"\n")
		puFile.WriteString(pu+"\n")
	}

}

func generatePrPuKeyPair() (string, string){
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)
	//fmt.Println(hexutil.Encode(privateKeyBytes)[2:]) // 0xfad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	//fmt.Println(address)
	return hexutil.Encode(privateKeyBytes)[2:], address
}