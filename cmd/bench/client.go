package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/smallnest/rpcx/client"
)

var (
	addr        = flag.String("addr", "tcp@localhost:8972", "snowflake service addressed")
	concurrency = flag.Int("c", 256, "concurrent clients")
	count       = flag.Int("n", 10000000, "fetched ids per client")
	batch       = flag.Int("b", 1, "batch number")
)

func main() {
	flag.Parse()

	addrs := strings.Split(*addr, ",")
	peers := make([]*client.KVPair, 0, len(addrs))
	for _, a := range addrs {
		peers = append(peers, &client.KVPair{Key: a})
	}

	var wg, latch sync.WaitGroup
	wg.Add(*concurrency)
	latch.Add(*concurrency)

	for i := 0; i < *concurrency; i++ {
		go func() {
			var id int64
			var ids []int64

			d := client.NewMultipleServersDiscovery(peers)
			xclient := client.NewXClient("snowflake", client.Failover, client.RoundRobin, d, client.DefaultOption)
			defer xclient.Close()

			//warmup
			xclient.Call(context.Background(), "Get", 0, &id)

			// start
			latch.Done()
			latch.Wait()

			n := *count / (*batch)
			for j := 0; j < n; j++ {
				if *batch <= 1 {
					err := xclient.Call(context.Background(), "Get", 0, &id)
					if err != nil {
						log.Fatalf("failed to call: %v", err)
					}
				} else {
					err := xclient.Call(context.Background(), "GetBatch", *batch, &ids)
					if err != nil {
						log.Fatalf("failed to call: %v", err)
					}
				}

			}

			wg.Done()
		}()
	}

	latch.Wait()
	start := time.Now()

	wg.Wait()
	dur := time.Since(start)

	fmt.Printf("total IDs: %d, duration: %v, id/s: %d", (*concurrency)*(*count), dur, int64(*concurrency)*int64(*count)*1000000/dur.Microseconds())
}
