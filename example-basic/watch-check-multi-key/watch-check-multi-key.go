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
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer etcdClient.Close()

	go func() {
		watcherClient, err := clientv3.New(clientv3.Config{
			Endpoints:   []string{"localhost:2379"},
			DialTimeout: 5 * time.Second,
		})
		if err != nil {
			log.Fatal(err)
		}
		defer watcherClient.Close()

		rch := watcherClient.Watch(context.Background(), "networking-rule", clientv3.WithPrefix())
		for wresp := range rch {
			for _, ev := range wresp.Events {
				fmt.Printf("Watcher (networking rule) - %s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			}
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

		rch := watcherClient.Watch(context.Background(), "networking-rule/group1", clientv3.WithPrefix())
		for wresp := range rch {
			for _, ev := range wresp.Events {
				fmt.Printf("Watcher (group)- %s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			}
		}
	}()

	time.Sleep(1 * time.Second)

	go func() {
		for i := 1; i <= 5; i++ {
			key := "networking-rule/group" + strconv.Itoa(i)
			etcdClient.Put(context.Background(), key, strconv.Itoa(i))
			time.Sleep(10 * time.Millisecond)
		}
	}()
	time.Sleep(1 * time.Second)

	go func() {
		for i := 1; i <= 5; i++ {
			key := "networking-rule/group1/host" + strconv.Itoa(i)
			etcdClient.Put(context.Background(), key, strconv.Itoa(i))
			time.Sleep(10 * time.Millisecond)
		}
	}()
	time.Sleep(1 * time.Second)

	// Get
	resp, respErr := etcdClient.Get(context.TODO(), "networking-rule/group", clientv3.WithPrefix())
	if respErr != nil {
		fmt.Println(respErr)
	}

	for _, kv := range resp.Kvs {
		fmt.Printf("Get - %q : %q\n", kv.Key, kv.Value)
	}

	// Get
	r, rErr := etcdClient.Get(context.TODO(), "networking-rule/group", clientv3.WithPrefix(), clientv3.WithCountOnly())
	if rErr != nil {
		fmt.Println(rErr)
	}
	fmt.Printf("%#v\n", r)
	fmt.Printf("Count: %#v\n", r.Count)

	var ch chan bool
	<-ch // Block forever
}
