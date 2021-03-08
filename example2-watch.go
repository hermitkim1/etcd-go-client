package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/client/v3"
	"log"
	"strconv"
	"time"
)


func main() {
	putterClient, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer putterClient.Close()

	go func() {
		for i := 0; i < 100; i++ {
			putterClient.Put(context.Background(), "foo", strconv.Itoa(i))
			time.Sleep(10 * time.Millisecond)
		}
	}()

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
				fmt.Printf("Watcher - %s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)

			}
		}
	}()

	var ch chan bool
	<- ch // blocks forever
}
