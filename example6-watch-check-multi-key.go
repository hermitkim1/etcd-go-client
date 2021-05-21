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
		watcherClient, err := clientv3.New(clientv3.Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
		})
		if err != nil {
			log.Fatal(err)
		}
		defer watcherClient.Close()


		rch := watcherClient.Watch(context.Background(), "foo", clientv3.WithPrefix())
		for wresp := range rch {
			for _, ev := range wresp.Events {
				fmt.Printf("Watcher - %s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)

			}
		}
	}()

	go func() {
		for i := 0; i < 5; i++ {
			key := "foo/" + strconv.Itoa(i)
			putterClient.Put(context.Background(), key, strconv.Itoa(i))
			time.Sleep(10 * time.Millisecond)
		}
	}()

	var ch chan bool
	<-ch // Block forever
}
