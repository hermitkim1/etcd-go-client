package main

import (
	"context"
	"flag"
	"fmt"
	"go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"log"
	"strconv"
	"time"
)


func main() {

	var myKey = "phoo"
	var name = flag.String("name", myKey, "give a name")
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

	// Create a sessions to aqcuire a lock
	session, _:= concurrency.NewSession(adderClient)
	defer session.Close()

	lock := concurrency.NewMutex(session, "/distributed-lock/")
	ctx := context.Background()

	// Set "phoo" with "50" initially
	adderClient.Put(ctx, myKey, strconv.Itoa(50))

	// Try to add value of "phoo" 100 times
	for i := 0; i < 100; i++ {
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
		time.Sleep(10 * time.Millisecond)

		// Release lock
		if err := lock.Unlock(ctx); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Released lock for ", *name)
	}

	fmt.Println("Done to add")

	var ch chan bool
	<- ch // Block forever
}
