package main

import (
	"context"
	"flag"
	"log"
	"strings"
	"time"

	"github.com/rpcxio/did"
	"github.com/smallnest/rpcx/client"
)

var (
	addr = flag.String("addr", "tcp@localhost:8972", "snowflake service addressed")
)

func main() {
	flag.Parse()

	addrs := strings.Split(*addr, ",")
	peers := make([]*client.KVPair, 0, len(addrs))
	for _, a := range addrs {
		peers = append(peers, &client.KVPair{Key: a})
	}

	d := client.NewMultipleServersDiscovery(peers)
	xclient := client.NewXClient("snowflake", client.Failover, client.RoundRobin, d, client.DefaultOption)
	defer xclient.Close()

	var id int64
	var ids []int64

	for {
		err := xclient.Call(context.Background(), "Get", 0, &id)
		if err != nil {
			log.Fatalf("failed to call: %v", err)
		}
		log.Printf("get a ID: %d (%s)", id, did.ID(id).Base58())
		time.Sleep(1e8)

		err = xclient.Call(context.Background(), "GetBatch", 10, &ids)
		if err != nil {
			log.Fatalf("failed to call: %v", err)
		}

		log.Printf("get a IDs: %v", ids)
		time.Sleep(1e8)
	}

}
