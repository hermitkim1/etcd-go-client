package main

import (
	"context"
	"etcd-go-client/configs"
	"fmt"
	"log"
	"path/filepath"
	"strconv"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func main() {

	// Load config
	configPath := filepath.Join("..", "..", "configs", "config.yaml")
	config, _ := configs.LoadConfig(configPath)

	var myKey = "phoo"

	// Subtractor Section
	subtractorClient, err := clientv3.New(clientv3.Config{
		Endpoints:   config.ETCD.Endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer subtractorClient.Close()

	fmt.Println("The subtractor is connected.")

	ctx := context.Background()

	for j := 0; j < 50; j++ {
		resp2, _ := subtractorClient.Get(ctx, myKey)
		num2, _ := strconv.Atoi(string(resp2.Kvs[0].Value))
		num2--
		fmt.Printf("Subtractordder: %v\n", num2)
		subtractorClient.Put(ctx, myKey, strconv.Itoa(num2))
		time.Sleep(10 * time.Millisecond)
	}

	fmt.Println("Done to subtract")
}
