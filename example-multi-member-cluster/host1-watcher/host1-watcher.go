package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/client/v3"
	"log"
	"time"
)


func main() {

	// Watcher Section
	watcherClient, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer watcherClient.Close()

	fmt.Println("The watcher is connected.")

	rch := watcherClient.Watch(context.Background(), "phoo")
	for wresp := range rch {
		for _, ev := range wresp.Events {
			fmt.Printf("Watch - %s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
		}
	}

	var ch chan bool
	<- ch // blocks forever
}
