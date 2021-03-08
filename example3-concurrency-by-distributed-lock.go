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

	var name = flag.String("name", "foo", "give a name")
	flag.Parse()

	// Watcher Section
	go func() {
		watcherClient, err := clientv3.New(clientv3.Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
		})
		if err != nil {
			log.Fatal(err)
		}
		defer watcherClient.Close()

		rch := watcherClient.Watch(context.Background(), "foo")
		for wresp := range rch {
			for _, ev := range wresp.Events {
				fmt.Printf("Watch - %s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			}
		}
	}()

	// Adder Section
	adderClient, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer adderClient.Close()

	// Create a sessions to aqcuire a lock
	s1, _:= concurrency.NewSession(adderClient)
	defer s1.Close()

	l1 := concurrency.NewMutex(s1, "/distributed-lock/")
	ctx1 := context.Background()

	adderClient.Put(context.Background(), "foo", strconv.Itoa(50))

	go func() {
		for i := 0; i < 100; i++ {
			// acquire lock (or wait to have it)
			if err := l1.Lock(ctx1); err != nil {
				log.Fatal(err)
			}
			fmt.Println("acquired lock for ", *name)

			resp1, _ := adderClient.Get(context.Background(), "foo")
			num1, _ := strconv.Atoi(string(resp1.Kvs[0].Value))
			num1++
			fmt.Printf("Adder: %v\n", num1)
			adderClient.Put(context.Background(), "foo", strconv.Itoa(num1))
			time.Sleep(10 * time.Millisecond)

			if err := l1.Unlock(ctx1); err != nil {
				log.Fatal(err)
			}
			fmt.Println("released lock for ", *name)
		}
	}()

	// Subtractor Section
	subtractorClient, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer subtractorClient.Close()

	// create a sessions to aqcuire a lock
	s2, _:= concurrency.NewSession(adderClient)
	defer s1.Close()

	l2 := concurrency.NewMutex(s2, "/distributed-lock/")
	ctx2 := context.Background()

	go func() {
		for j := 0; j < 100; j++ {
			// acquire lock (or wait to have it)
			if err := l2.Lock(ctx2); err != nil {
				log.Fatal(err)
			}
			fmt.Println("acquired lock for ", *name)

			resp2, _ := subtractorClient.Get(context.Background(), "foo")
			num2, _ := strconv.Atoi(string(resp2.Kvs[0].Value))
			num2--
			fmt.Printf("Subtractordder: %v\n", num2)
			subtractorClient.Put(context.Background(), "foo", strconv.Itoa(num2))
			time.Sleep(10 * time.Millisecond)

			if err := l2.Unlock(ctx2); err != nil {
				log.Fatal(err)
			}
			fmt.Println("released lock for ", *name)
		}
	}()

	var ch chan bool
	<- ch // blocks forever
}
