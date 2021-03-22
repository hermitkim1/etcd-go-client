package main

import (
	"context"
	"flag"
	"fmt"
	"go.etcd.io/etcd/client/v3"
	"log"
	"strconv"
	"time"
)

func main() {

	var myKey = "phoo"
	flag.Parse()

	// Adder Section
	adderClient, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer adderClient.Close()

	fmt.Println("Adder is connected.")

	requestTimeout := 10 * time.Second


	// Set "phoo" with "50" initially
	adderClient.Put(context.Background(), myKey, strconv.Itoa(50))

	// Try to add(+1) value of "phoo" 50 times
	for i := 0; i < 500; i++ {
		// Try compare-and-swap until succeeded
		for {
			resp, _ := adderClient.Get(context.Background(), myKey)
			value := string(resp.Kvs[0].Value)
			num, _ := strconv.Atoi(value)
			num++
			fmt.Printf("[Adder] Value1(%v), num(%v)\n", value, num)

			// Compare-and-Swap (CAS)
			ctx, _ := context.WithTimeout(context.TODO(), requestTimeout)
			txResp, err2 := adderClient.Txn(ctx).
				If(clientv3.Compare(clientv3.Value(myKey), "=", value)).
				Then(clientv3.OpPut(myKey, strconv.Itoa(num))).
				Else(clientv3.OpGet(myKey)).
				Commit()

			if err2 != nil {
				fmt.Printf("[Adder] i(%v) - Err1: %v\n", i, err2)
				//cancel1()

			}
			fmt.Printf("[Adder] i(%v) - txResp: %v\n", i, txResp)

			if txResp.Succeeded {
				break
			}

			//time.Sleep(1 * time.Millisecond)
		}
		//time.Sleep(1 * time.Millisecond)
	}

	fmt.Println("Done to add")
}
