// Run multiple Unary RPC

package main

import (
	"time"

	pb "github.com/hmarui66/grpc-sample/proto"
	"google.golang.org/grpc"
)

// unaryMulti run Unary RPC multiple times with a 3-second wait
func unaryMulti(client pb.RouteGuideClient) {
	for i := 0; i < 10; i++ {
		// Looking for a valid feature
		printFeature(client, &pb.Point{Latitude: 409146138, Longitude: -746188906})
		time.Sleep(3 * time.Second)
	}
}

// unaryMultiNonReuseCli run Unary RPC multiple times with a 3-second wait without reusing clients
func unaryMultiNonReuseCli(conn *grpc.ClientConn) {
	for i := 0; i < 10; i++ {
		client := pb.NewRouteGuideClient(conn)
		// Looking for a valid feature
		printFeature(client, &pb.Point{Latitude: 409146138, Longitude: -746188906})
		time.Sleep(3 * time.Second)
	}
}

// unaryMultiNonReuseConn run Unary RPC multiple times with a 3-second wait without reusing connection
func unaryMultiNonReuseConn(opts []grpc.DialOption) {
	for i := 0; i < 10; i++ {
		func() {
			conn := conn(opts)
			defer closerClose(conn)
			client := pb.NewRouteGuideClient(conn)
			// Looking for a valid feature
			printFeature(client, &pb.Point{Latitude: 409146138, Longitude: -746188906})
			time.Sleep(3 * time.Second)
		}()
	}
}
