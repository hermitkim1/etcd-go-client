package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

func main() {

	// var name = flag.String("name", "foo", "Give a name")
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

		rch := watcherClient.Watch(context.Background(), "root", clientv3.WithPrefix())
		for wresp := range rch {
			fmt.Printf("%#v\n", wresp)
			for _, ev := range wresp.Events {
				fmt.Printf("Watch - %s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			}
		}
	}()

	time.Sleep(3 * time.Second)
	// Adder Section
	updaterClient, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer updaterClient.Close()

	// Create a sessions to aqcuire a lock
	s1, _ := concurrency.NewSession(updaterClient)
	defer s1.Close()

	// l1 := concurrency.NewMutex(s1, "/distributed-lock/")
	// ctx1 := context.Background()

	updaterClient.Put(context.Background(), "root", "I'm root")
	time.Sleep(1 * time.Second)
	updaterClient.Put(context.Background(), "root/d1-key1", "I'm d1-key1")
	time.Sleep(1 * time.Second)
	updaterClient.Put(context.Background(), "root/d1-key1/d2-key1", "I'm d2-key1")
	time.Sleep(1 * time.Second)
	updaterClient.Put(context.Background(), "root/d1-key1/d2-key2", "I'm d2-key2")

	var ch chan bool
	<-ch // Block forever
}
