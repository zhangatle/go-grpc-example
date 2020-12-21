package main

import (
	"context"
	"go-grpc-example/pkg/gtls"
	pb "go-grpc-example/proto"
	"google.golang.org/grpc"
	"log"
)
const PORT = "9001"

func main() {
	tlsClient := gtls.Client{
		ServerName: "grpc.pro",
		CertFile : "../../conf/ca/intermediate/certs/grpc.pro.cert.pem",
		KeyFile : "../../conf/ca/intermediate/private/grpc.pro.key.pem",
		CaFile : "../../conf/ca/intermediate/certs/ca-chain.cert.pem",
	}
	c, err := tlsClient.GetTLSCredentials()
	if err != nil {
		log.Fatalf("tlsClient.GetTLSClientials err: %v", err)
	}

	conn, err := grpc.Dial(":"+PORT, grpc.WithTransportCredentials(c))
	if err != nil {
		log.Fatalf("grpc.Dial err: %v", err)
	}
	defer conn.Close()

	client := pb.NewSearchServiceClient(conn)
	resp, err := client.Search(context.Background(), &pb.SearchRequest{
		Request: "gRPC",
	})
	if err != nil {
		log.Fatalf("client.Search err: %v", err)
	}

	log.Printf("resp: %s", resp.GetResponse())
}
