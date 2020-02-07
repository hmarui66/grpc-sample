package main

import (
	"time"

	pb "github.com/hmarui66/grpc-sample/proto"
)

// keep run Unary RPC twice with 1 minute wait
func keep(client pb.RouteGuideClient) {
	printFeature(client, &pb.Point{Latitude: 409146138, Longitude: -746188906})

	time.Sleep(60 * time.Second)

	printFeature(client, &pb.Point{Latitude: 409146138, Longitude: -746188906})
}

// keep run a Unary RPC after 1 minute wait
func keepWithoutFirstCall(client pb.RouteGuideClient) {
	time.Sleep(60 * time.Second)

	printFeature(client, &pb.Point{Latitude: 409146138, Longitude: -746188906})
}
