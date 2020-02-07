package main

import (
	"math/rand"
	"sync"
	"time"

	pb "github.com/hmarui66/grpc-sample/proto"
	"google.golang.org/grpc"
)

// streamImage run Stream RPC to upload image
func streamImage(client pb.RouteGuideClient) {
	runRecordImage(client)
}

// streamImage run Stream RPC to upload image
func manyConnStreamImage(opts []grpc.DialOption) {
	var wg sync.WaitGroup

	for i := 0; i < 500; i++ {
		wg.Add(1)
		go func() {
			conn := conn(opts)
			defer closerClose(conn)
			client := pb.NewRouteGuideClient(conn)

			for j := 0; j < 100; j++ {
				time.Sleep(time.Duration(rand.Int31n(5)) * time.Second)
				runRecordImage(client)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
