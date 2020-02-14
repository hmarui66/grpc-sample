package main

import (
	"context"
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	pb "github.com/hmarui66/grpc-sample/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/testdata"
)

var (
	tls                = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	caFile             = flag.String("ca_file", "", "The file containing the CA root cert file")
	serverAddr         = flag.String("server_addr", ":8080", "The server address in the format of host:port")
	serverHostOverride = flag.String("server_host_override", "x.test.youtube.com", "The server name use to verify the hostname returned by TLS handshake")
)

func printFeature(client pb.RouteGuideClient, point *pb.Point) {
	log.Printf("Getting feature for point (%d, %d)", point.Latitude, point.Longitude)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	feature, err := client.GetFeature(ctx, point)
	if err != nil {
		log.Printf("%v.GetFeatures(_) = _, %v: ", client, err)
		return
	}
	log.Println(feature)
}

func printFeatures(client pb.RouteGuideClient, rect *pb.Rectangle) {
	log.Printf("Looking for features within %v", rect)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := client.ListFeatures(ctx, rect)
	if err != nil {
		log.Printf("%v.ListFeatures(_) = _, %v", client, err)
		return
	}
	for {
		feature, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("%v.ListFeatures(_) = _, %v", client, err)
			return
		}
		log.Println(feature)
	}
}

func runRecordImage(client pb.RouteGuideClient) {
	f, err := os.Open(filepath.Join("assets", "sample.jpg"))
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}
	defer closerClose(f)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := client.RecordLocationImage(ctx)
	if err != nil {
		log.Printf("%v.RecordRoute(_) = _, %v", client, err)
		return
	}

	var b [4096 * 1000]byte
	for {
		n, err := f.Read(b[:])
		if err != nil {
			if err != io.EOF {
				log.Fatalf("failed to read image data: %v", err)
			}
			break
		}
		err = stream.Send(&pb.Image{
			Data: b[:n],
		})
		if err != nil {
			log.Printf("failed to send a chunked data: %v", err)
		}
	}
	reply, err := stream.CloseAndRecv()
	if err != nil {
		log.Printf("%v.CloseAndRecv() got error %v, want %v", stream, err, nil)
		return
	}
	log.Printf("Image size: %v", reply)
}

func closerClose(c io.Closer) {
	if err := c.Close(); err != nil {
		log.Fatalf("failed to close: %v", err)
	}
}

func conn(opts []grpc.DialOption) *grpc.ClientConn {
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	return conn
}

func main() {
	flag.Parse()
	var opts []grpc.DialOption
	if *tls {
		if *caFile == "" {
			*caFile = testdata.Path("ca.pem")
		}
		creds, err := credentials.NewClientTLSFromFile(*caFile, *serverHostOverride)
		if err != nil {
			log.Fatalf("failed to create TLS credentials %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	opts = append(opts, grpc.WithBlock())
	conn := conn(opts)
	defer closerClose(conn)
	client := pb.NewRouteGuideClient(conn)

	switch flag.Arg(0) {
	case "unary":
		unaryMulti(client)
	case "unary-non-reuse-cli":
		unaryMultiNonReuseCli(conn)
	case "unary-non-reuse-conn":
		closerClose(conn)
		unaryMultiNonReuseConn(opts)
	case "stream":
		streamMulti(client)
	case "stream-non-reuse-cli":
		streamMultiNonReuseCli(conn)
	case "stream-non-reuse-conn":
		closerClose(conn)
		streamMultiNonReuseConn(opts)
	case "keep":
		keep(client)
	case "keep-without-first-call":
		keepWithoutFirstCall(client)
	case "many-conn":
		closerClose(conn)
		manyConn(opts)
	case "many-conn-stream":
		closerClose(conn)
		manyConnStream(opts)
	case "stream-image":
		closerClose(conn)
		streamImage(client)
	case "many-conn-stream-image":
		closerClose(conn)
		manyConnStreamImage(opts)
	default:
		log.Fatalf("invalid command args")
	}
}
