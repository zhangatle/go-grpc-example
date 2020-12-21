package main

import (
	"context"
	"go-grpc-example/pkg/gtls"
	pb "go-grpc-example/proto"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
)

type SearchService struct{
	pb.UnimplementedSearchServiceServer
}

func (s *SearchService) Search(ctx context.Context, r *pb.SearchRequest) (*pb.SearchResponse, error) {
	return &pb.SearchResponse{Response: r.GetRequest() + " Server"}, nil
}

const PORT = "9001"

func main() {
	certFile := "../../conf/ca/intermediate/certs/grpc.pro.cert.pem"
	keyFile := "../../conf/ca/intermediate/private/grpc.pro.key.pem"
	caFile := "../../conf/ca/intermediate/certs/ca-chain.cert.pem"
	tlsServer := gtls.Server{
		CertFile: certFile,
		KeyFile:  keyFile,
		CaFile: caFile,
	}
	c, err := tlsServer.GetTLSCredentials()
	if err != nil {
		log.Fatalf("tlsServer.GetTLSCredentials err: %v", err)
	}
	//mux := GetHTTPServeMux()
	server := grpc.NewServer(grpc.Creds(c))
	pb.RegisterSearchServiceServer(server, &SearchService{})

	//_ = http.ListenAndServeTLS(":"+PORT,
	//	certFile,
	//	keyFile,
	//	http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
	//			server.ServeHTTP(w, r)
	//		} else {
	//			//mux.ServeHTTP(w, r)
	//		}
	//		return
	//	}),
	//)
	lis, err := net.Listen("tcp", ":"+PORT)
	if err != nil {
		log.Fatalf("net.Listen err: %v", err)
	}
	_ = server.Serve(lis)
}

func GetHTTPServeMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("grpc.pro"))
	})
	return mux
}
