package api

import (
	"log"
	"net"
	"strconv"

	grpcmw "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/happilymarrieddad/old-world/api3/internal/api/auth"
	interceptors "github.com/happilymarrieddad/old-world/api3/internal/api/interceptors"
	"github.com/happilymarrieddad/old-world/api3/internal/repos"
	pbauth "github.com/happilymarrieddad/old-world/api3/pb/proto/auth"
	"google.golang.org/grpc"
)

func Run(gr repos.GlobalRepo) {
	port := int(50051)
	lis, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		log.Fatal("failed to listen: " + err.Error())
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpcmw.ChainUnaryServer(
				interceptors.GlobalRepoInjector(gr),
			),
		),
	)

	pbauth.RegisterV1AuthServer(s, auth.InitRoutes())

	log.Printf("Server listening on port %d\n", port)
	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
