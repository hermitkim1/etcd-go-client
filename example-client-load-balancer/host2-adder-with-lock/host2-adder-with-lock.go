package main

import (
	"context"
	"etcd-go-client/configs"
	"flag"
	"fmt"
	"go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
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
	var name = flag.String("name", myKey, "give a name")
	flag.Parse()

	// Adder Section
	adderClient, err := clientv3.New(clientv3.Config{
		Endpoints:   config.ETCD.Endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer adderClient.Close()

	fmt.Println("Adder is connected.")

	// Create a sessions to aqcuire a lock
	session, _:= concurrency.NewSession(adderClient)
	defer session.Close()

	lock := concurrency.NewMutex(session, "/distributed-lock/")
	ctx := context.Background()

	// Set "phoo" with "50" initially
	adderClient.Put(ctx, myKey, strconv.Itoa(50))

	// Try to add value of "phoo" 100 times
	for i := 0; i < 50; i++ {
		// Acquire lock (or wait to have it)
		if err := lock.Lock(ctx); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Acquired lock for ", *name)

		resp1, _ := adderClient.Get(ctx, myKey)
		num1, _ := strconv.Atoi(string(resp1.Kvs[0].Value))
		num1++
		fmt.Printf("Adder: %v\n", num1)
		adderClient.Put(ctx, myKey, strconv.Itoa(num1))
		//time.Sleep(1 * time.Millisecond)

		// Release lock
		if err := lock.Unlock(ctx); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Released lock for ", *name)
	}

	fmt.Println("Done to add")
}
