package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	"sync"
	"time"
)

func main(){

	client, err := ethclient.Dial("https://mainnet.infura.io/v3/6967c52c1d214d0a877ff44e7759867d")
	if err != nil {
		log.Fatal(err)
	}
    GetRequiredBlocks(client)

}


func GetRequiredBlocks(client *ethclient.Client) {

	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	header.Hash()
	//var index = 0
    //requiredBlockArray := make([]int64,blockQuantity+2)
	latestBlockNumber := header.Number.Int64()
    RunInConcurrency(client,latestBlockNumber,100)

}
func RunInConcurrency(client *ethclient.Client,blockNumber int64,totalBlocksRequired int64){
	var wg sync.WaitGroup

	startTime := time.Now()
	 var newBlockNumber   = blockNumber
	 requiredLoops :=  totalBlocksRequired/10
	 for i := int64(0) ; i < requiredLoops ; i++{
		 wg.Add(1)
		 go FetchBlocks(client,newBlockNumber,&wg)
	 	   newBlockNumber = newBlockNumber - 11
 	 }
	wg.Wait()
	fmt.Println(time.Now().Sub(startTime))

}

func FetchBlocks(client *ethclient.Client,blockNumber int64,wg *sync.WaitGroup)  []string{
	defer wg.Done()
	var index = 0
	result :=  Result{}
	requiredBlockArray := make([]string,11)
	for i := blockNumber ; i >= blockNumber - 10 ; i-- {
		block, err := client.BlockByNumber(context.Background(), big.NewInt(i))
		if err != nil {
			log.Fatal(err)

		}
		marshall ,_ := block.Bloom().MarshalText()
		data := BlockInfo{BlockNumber: i,
			Hash: block.Hash().String(),Logs: string(marshall),
		}
		result.BlockData = append(result.BlockData,data)
		requiredBlockArray[index] = block.Hash().String()
		index ++
	}
	fmt.Println(result)
	return requiredBlockArray
}