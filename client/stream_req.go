package main

import (
	"time"

	pb "github.com/hmarui66/grpc-sample/proto"
	"google.golang.org/grpc"
)

// streamMulti run Stream RPC multiple times with a 3-second wait
func streamMulti(client pb.RouteGuideClient) {
	for i := 0; i < 10; i++ {
		// Looking for a valid feature
		printFeatures(client, &pb.Rectangle{
			Lo: &pb.Point{Latitude: 400000000, Longitude: -750000000},
			Hi: &pb.Point{Latitude: 420000000, Longitude: -730000000},
		})
		time.Sleep(3 * time.Second)
	}
}

// streamMultiNonReuseCli run Stream RPC multiple times with a 3-second wait without reusing clients
func streamMultiNonReuseCli(conn *grpc.ClientConn) {
	for i := 0; i < 10; i++ {
		client := pb.NewRouteGuideClient(conn)
		// Looking for a valid feature
		printFeatures(client, &pb.Rectangle{
			Lo: &pb.Point{Latitude: 400000000, Longitude: -750000000},
			Hi: &pb.Point{Latitude: 420000000, Longitude: -730000000},
		})
		time.Sleep(3 * time.Second)
	}
}

// streamMultiNonReuseConn run Stream RPC multiple times with a 3-second wait without reusing connection
func streamMultiNonReuseConn(opts []grpc.DialOption) {
	for i := 0; i < 10; i++ {
		func() {
			conn := conn(opts)
			defer closerClose(conn)
			client := pb.NewRouteGuideClient(conn)
			// Looking for a valid feature
			printFeatures(client, &pb.Rectangle{
				Lo: &pb.Point{Latitude: 400000000, Longitude: -750000000},
				Hi: &pb.Point{Latitude: 420000000, Longitude: -730000000},
			})
			time.Sleep(3 * time.Second)
		}()
	}
}
