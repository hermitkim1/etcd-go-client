package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

func mockController(wg *sync.WaitGroup, name string, keyToTest string, keyToUpdate string) {
	defer wg.Done()

	fmt.Printf("##### Start ---------- mockController (%s)\n", name)

	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer etcdClient.Close()

	// create a sessions to aqcuire a lock
	s, _ := concurrency.NewSession(etcdClient)
	defer s.Close()

	rch := etcdClient.Watch(context.Background(), keyToTest)
	for wresp := range rch {
		for _, ev := range wresp.Events {

			fmt.Printf("\n[%s] Watch - %s %q : %q\n", name, ev.Type, ev.Kv.Key, ev.Kv.Value)
			keyPrefix := fmt.Sprintf("lease/%s-%d", keyToUpdate, wresp.Header.GetRevision())
			// fmt.Printf("%#v\n", keyPrefix)

			// Self-assign a task by Compare-and-Swap (CAS) and Lease
			lease := clientv3.NewLease(etcdClient)
			ttl := int64(15)
			grantResp, err := lease.Grant(context.TODO(), ttl)
			if err != nil {
				fmt.Printf("[%s] 'lease.Grant' error: %#v\n", name, err)
			}

			messageToCheck := fmt.Sprintf("Vanish in %d sec", ttl)
			txResp, err2 := etcdClient.Txn(context.TODO()).
				If(clientv3.Compare(clientv3.Value(keyPrefix), "=", messageToCheck)).
				Then(clientv3.OpGet(keyPrefix)).
				Else(clientv3.OpPut(keyPrefix, messageToCheck, clientv3.WithLease(grantResp.ID))).
				Commit()

			if err2 != nil {
				fmt.Printf("[%s] transaction error: %#v\n", name, err2)
			}

			// fmt.Printf("[%s] txResp: %v\n", name, txResp)
			needToHandle := !txResp.Succeeded

			if needToHandle {

				l := concurrency.NewMutex(s, keyToUpdate)

				ctx := context.TODO()
				// Acquire lock (or wait to have it)
				err := l.Lock(ctx)
				if err != nil {
					fmt.Printf("[%s] '%s' is locked\n", name, keyToUpdate)
					log.Println(err)
				} else {
					// fmt.Printf("[%s] Acquired lock for %s\n", name, keyPrefix)
					// fmt.Printf("[%s] Do some work in %s\n", name, keyPrefix)

					// Do value+1 and put the value
					resp1, _ := etcdClient.Get(context.Background(), keyToUpdate)
					num1, _ := strconv.Atoi(string(resp1.Kvs[0].Value))
					num1++
					fmt.Printf("[%s] Adder: %v\n", name, num1)
					etcdClient.Put(context.Background(), keyToUpdate, strconv.Itoa(num1))
					rand.Seed(time.Now().UnixNano())
					t := rand.Int63n(401) + 100
					time.Sleep(time.Duration(t) * time.Millisecond)

					if err := l.Unlock(ctx); err != nil {
						log.Fatal(err)
					}
					// fmt.Printf("[%s] Released lock for %s\n", name, keyPrefix)
				}
			}
		}
	}
	fmt.Printf("##### End ---------- mockController (%s)\n", name)
}

func main() {

	var keyToWatch = "key-watched"
	var keyToUpdate = "key-updated"

	// Wait for multiple goroutines to complete
	var wg sync.WaitGroup

	// (a controller) Watch and do something
	wg.Add(3)
	go mockController(&wg, "controller 1", keyToWatch, keyToUpdate)
	go mockController(&wg, "controller 2", keyToWatch, keyToUpdate)
	go mockController(&wg, "controller 3", keyToWatch, keyToUpdate)

	// (an agent) Update
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer etcdClient.Close()

	// Initialize a value of "key-updated"
	etcdClient.Put(context.Background(), keyToUpdate, "0")

	// Wait 5 seconds until goroutines are ready
	time.Sleep(3 * time.Second)

	// Set a value of "key-watched"
	for i := 1; i <= 30; i++ {
		etcdClient.Put(context.Background(), keyToWatch, strconv.Itoa(i))
	}

	// Wait 3 seconds until goroutines are ready
	time.Sleep(10 * time.Second)

	// Results
	resp1, _ := etcdClient.Get(context.Background(), keyToWatch)
	fmt.Printf("Value of 'key-watched': %s\n", string(resp1.Kvs[0].Value))

	resp2, _ := etcdClient.Get(context.Background(), keyToUpdate)
	fmt.Printf("Value of 'key-updated': %s\n", string(resp2.Kvs[0].Value))

	wg.Wait()
}
