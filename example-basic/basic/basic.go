package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/client/v3"
	"time"
)

func main() {
	// Expect dial time-out on ipv4 blackhole
	_, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"http://localhost:2373"},
		DialTimeout: 2 * time.Second,
	})



	// etcd clientv3 >= v3.2.10, grpc/grpc-go >= v1.7.3
	if err == context.DeadlineExceeded {
		// Handle errors
	}

	cli2, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		// Handle error!
	}
	defer cli2.Close()

	//ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	_, err = cli2.Put(cli2.Ctx(), "sample_key", "sample_value")
	if err != nil {
		// Handle error!
	}

	resp, err2 := cli2.Get(cli2.Ctx(), "sample_key")
	if err2 != nil {
		// Handle error!
	}
	// Use the response
	fmt.Println(resp.Header)
	fmt.Println(resp.Kvs)

}
