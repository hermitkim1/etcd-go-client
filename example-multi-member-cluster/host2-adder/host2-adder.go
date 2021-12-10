package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func main() {

	var myKey = "phoo"

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

	ctx := context.Background()

	// Set "phoo" with "50" initially
	adderClient.Put(ctx, myKey, strconv.Itoa(50))

	// Try to add value of "phoo" 100 times
	for i := 0; i < 50; i++ {
		resp1, _ := adderClient.Get(ctx, myKey)
		num1, _ := strconv.Atoi(string(resp1.Kvs[0].Value))
		num1++
		fmt.Printf("Adder: %v\n", num1)
		adderClient.Put(ctx, myKey, strconv.Itoa(num1))
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Done to add")
}
