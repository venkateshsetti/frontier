package blockQuery

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"
	"math/big"
	"sync"
	"time"
)

// Manager  parameters
type Manager struct {
	log *zap.SugaredLogger
}

// NewManager Assigning the parameters to Manager instance
func NewManager(logger *zap.SugaredLogger) *Manager {
	mgr := &Manager{log: logger}
	return mgr
}


func (m *Manager) GetRequiredBlocks() Response {
	var response Response
	client, err := ethclient.Dial("https://mainnet.infura.io/v3/8a5b0f66345940ccb7f73113d13b34fb")  //creating the ethereum client by using project ID
	if err != nil {
		m.log.Fatal(err)
	}
	header, err := client.HeaderByNumber(context.Background(), nil) //Getting the latest BlockNumber
	if err != nil {
		m.log.Fatal(err)
	}
	latestBlockNumber := header.Number.Int64() //converting it to int type
	m.log.Infof("Latest Block Number %d",latestBlockNumber)
	var blockData []BlockInfo
	res := m.RunInConcurrency(client,latestBlockNumber,10000)  //Passing the latest block number and no of records need to be fetched
	for _,data  := range res {
        block := data[latestBlockNumber]  //Passing the block number to get the response associated to that map key
        blockData = append(blockData,block.Result...)  //appending  the slice of block data
        latestBlockNumber = latestBlockNumber - 11 //updating the key value to get the another map data
	}
	response.Message = "Success"
	response.Result = blockData
	return response
}



func (m *Manager) RunInConcurrency(client *ethclient.Client,blockNumber int64,totalBlocksRequired int64) []map[int64]Response{
	var wg sync.WaitGroup  //synchronize lets the main go routine to wait till it complete the all remaining go routines

	startTime := time.Now()
	sortResponse := make([]map[int64]Response,1001)  //creating the map with capacity 1001
	requiredLoops :=  totalBlocksRequired/10   //decides how many go routines required for the batch of records
	/*
	  Following Loop will run the required go routines and also updates block number following divide and conquer approach
	  and also storing the every go routine response into the slice
	 */
	for i,j := int64(0), blockNumber ; i < requiredLoops ; i,j = i+1,j-11 {
		wg.Add(1)
		j := j
		i := i
		go func(){
			sortResponse[i] = m.FetchBlocks(client,j)
			wg.Done()
		}()
	}
	wg.Wait()
	m.log.Infof("Total Time Taken %s",time.Now().Sub(startTime)) //Total time taken for all go routines to  fetch the block records
	return sortResponse

}

/*
   Following Method calls ethereum BlockByNumber by new block number  inside a loop
   and fetches the data and parsing the data into BlockInfo struct
   and appending all the result fethced from requests
   assigning the appended data to the starting block number as key and the appended data as value
 */
func (m *Manager) FetchBlocks(client *ethclient.Client,blockNumber int64) map[int64]Response{
	result :=  Response{}
	mapBlockArray := make(map[int64]Response)
	for i := blockNumber ; i >= blockNumber - 10 ; i-- {
		block, err := client.BlockByNumber(context.Background(), big.NewInt(i))
		if err != nil {
			m.log.Fatal(err)
		}
		marshall ,_ := block.Bloom().MarshalText()
		data := BlockInfo{BlockNumber: fmt.Sprintf("%X",i),
			Hash: block.Hash().String(),Logs: string(marshall),
		}
		result.Result = append(result.Result,data)
	}
	mapBlockArray[blockNumber] = result
	return mapBlockArray
}