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
	var name = flag.String("name", myKey, "Give a name")
	flag.Parse()

	// Subtractor Section
	subtractorClient, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer subtractorClient.Close()

	fmt.Println("The subtractor is connected.")

	// Create a sessions to aqcuire a lock
	session, _:= concurrency.NewSession(subtractorClient)
	defer session.Close()

	lock := concurrency.NewMutex(session, "/distributed-lock/")
	ctx := context.Background()

	go func() {
		for j := 0; j < 100; j++ {
			// Acquire lock (or wait to have it)
			if err := lock.Lock(ctx); err != nil {
				log.Fatal(err)
			}
			fmt.Println("Acquired lock for ", *name)

			resp2, _ := subtractorClient.Get(context.Background(), myKey)
			num2, _ := strconv.Atoi(string(resp2.Kvs[0].Value))
			num2--
			fmt.Printf("Subtractordder: %v\n", num2)
			subtractorClient.Put(context.Background(), myKey, strconv.Itoa(num2))
			time.Sleep(10 * time.Millisecond)

			// Release lock
			if err := lock.Unlock(ctx); err != nil {
				log.Fatal(err)
			}
			fmt.Println("Released lock for ", *name)
		}
	}()

	var ch chan bool
	<- ch // Block forever
}
