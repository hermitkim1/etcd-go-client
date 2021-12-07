package main

import (
	"context"
	"etcd-go-client/configs"
	"flag"
	"fmt"
	"go.etcd.io/etcd/client/v3"
	"log"
	"path/filepath"
	"strconv"
	"time"
)

func main() {

	// Load config
	configPath := filepath.Join("..", "..", "configs", "config.yaml")
	config, _ := configs.LoadConfig(configPath)

	var myKey = "phoo"
	flag.Parse()

	// Subtractor Section
	subtractorClient, err := clientv3.New(clientv3.Config{
		Endpoints:   config.ETCD.Endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer subtractorClient.Close()

	fmt.Println("The subtractor is connected.")

	requestTimeout := 10 * time.Second

	for j := 0; j < 50; j++ {
		// Try compare-and-swap until succeeded
		for {
			resp, _ := subtractorClient.Get(context.Background(), myKey)
			value := string(resp.Kvs[0].Value)
			num, _ := strconv.Atoi(value)
			num--
			fmt.Printf("[Subtractor] Value2(%v), num(%v)\n", value, num)

			// Compare-and-Swap (CAS)
			ctx, _ := context.WithTimeout(context.TODO(), requestTimeout)
			txResp2, err2 := subtractorClient.Txn(ctx).
				If(clientv3.Compare(clientv3.Value(myKey), "=", value)).
				Then(clientv3.OpPut(myKey, strconv.Itoa(num))).
				Else(clientv3.OpGet(myKey)).
				Commit()

			if err2 != nil {
				fmt.Printf("[Subtractor] j(%v) - Err2: %v\n", j, err2)
			}

			fmt.Printf("[Subtractor] j(%v) - txResp2: %v\n", j, txResp2)

			if txResp2.Succeeded {
				break
			}
			//time.Sleep(1 * time.Millisecond)
		}
		//time.Sleep(1 * time.Millisecond)
	}

	fmt.Println("Done to subtract")
}
