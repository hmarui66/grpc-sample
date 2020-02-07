package main

import (
	"math/rand"
	"sync"
	"time"

	pb "github.com/hmarui66/grpc-sample/proto"
	"google.golang.org/grpc"
)

// manyConn generates many connections and run Unary RPC on each connection
func manyConn(opts []grpc.DialOption) {
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			conn := conn(opts)
			defer closerClose(conn)
			client := pb.NewRouteGuideClient(conn)

			for j := 0; j < 100; j++ {
				time.Sleep(time.Duration(rand.Int31n(5)) * time.Second)
				printFeature(client, &pb.Point{Latitude: 409146138, Longitude: -746188906})
			}
			wg.Done()
		}()
	}

	wg.Wait()
}

// manyConnStream generates many connections and run Stream RPC on each connection
func manyConnStream(opts []grpc.DialOption) {
	var wg sync.WaitGroup

	for i := 0; i < 500; i++ {
		wg.Add(1)
		go func() {
			conn := conn(opts)
			defer closerClose(conn)
			client := pb.NewRouteGuideClient(conn)

			for j := 0; j < 100; j++ {
				time.Sleep(time.Duration(rand.Int31n(5)) * time.Second)
				printFeatures(client, &pb.Rectangle{
					Lo: &pb.Point{Latitude: 400000000, Longitude: -750000000},
					Hi: &pb.Point{Latitude: 420000000, Longitude: -730000000},
				})
			}
			wg.Done()
		}()
	}

	wg.Wait()
}
