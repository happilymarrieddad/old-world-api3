package api

import (
	"log"
	"net"
	"strconv"

	grpcmw "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/happilymarrieddad/old-world/api3/internal/api/auth"
	interceptors "github.com/happilymarrieddad/old-world/api3/internal/api/interceptors"
	v1armytypes "github.com/happilymarrieddad/old-world/api3/internal/api/v1/armytypes"
	v1compositiontypes "github.com/happilymarrieddad/old-world/api3/internal/api/v1/compositiontypes"
	v1games "github.com/happilymarrieddad/old-world/api3/internal/api/v1/games"
	v1itemtypes "github.com/happilymarrieddad/old-world/api3/internal/api/v1/itemtypes"
	v1optiontypes "github.com/happilymarrieddad/old-world/api3/internal/api/v1/optiontypes"
	v1statistics "github.com/happilymarrieddad/old-world/api3/internal/api/v1/statistics"
	v1trooptypes "github.com/happilymarrieddad/old-world/api3/internal/api/v1/trooptypes"
	v1unittypes "github.com/happilymarrieddad/old-world/api3/internal/api/v1/unittypes"
	v1userarmies "github.com/happilymarrieddad/old-world/api3/internal/api/v1/userarmies"
	v1userunits "github.com/happilymarrieddad/old-world/api3/internal/api/v1/userunits"
	"github.com/happilymarrieddad/old-world/api3/internal/repos"
	pbarmytypes "github.com/happilymarrieddad/old-world/api3/pb/proto/armytypes"
	pbauth "github.com/happilymarrieddad/old-world/api3/pb/proto/auth"
	pbcompositiontypes "github.com/happilymarrieddad/old-world/api3/pb/proto/compositiontypes"
	pbgames "github.com/happilymarrieddad/old-world/api3/pb/proto/games"
	pbitemtypes "github.com/happilymarrieddad/old-world/api3/pb/proto/itemtypes"
	pboptiontypes "github.com/happilymarrieddad/old-world/api3/pb/proto/optiontypes"
	pbstatistics "github.com/happilymarrieddad/old-world/api3/pb/proto/statistics"
	pbtrooptypes "github.com/happilymarrieddad/old-world/api3/pb/proto/trooptypes"
	pbunittypes "github.com/happilymarrieddad/old-world/api3/pb/proto/unittypes"
	pbuserarmies "github.com/happilymarrieddad/old-world/api3/pb/proto/userarmies"
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

	pbauth.RegisterAuthServer(s, auth.InitRoutes())

	// V1
	pbgames.RegisterV1GamesServer(s, v1games.InitRoutes())
	pbarmytypes.RegisterV1ArmyTypesServer(s, v1armytypes.InitRoutes())
	pbcompositiontypes.RegisterV1CompositionTypesServer(s, v1compositiontypes.InitRoutes())
	pbitemtypes.RegisterV1ItemTypesServer(s, v1itemtypes.InitRoutes())
	pboptiontypes.RegisterV1OptionTypesServer(s, v1optiontypes.InitRoutes())
	pbstatistics.RegisterV1StatisticsServer(s, v1statistics.InitRoutes())
	pbtrooptypes.RegisterV1TroopTypesServer(s, v1trooptypes.InitRoutes())
	pbunittypes.RegisterV1UnitTypesServer(s, v1unittypes.InitRoutes())
	pbuserarmies.RegisterV1UserArmiesServer(s, v1userarmies.InitRoutes())
	pbuserarmies.RegisterV1UserArmyUnitsServer(s, v1userunits.InitRoutes())

	log.Printf("Server listening on port %d\n", port)
	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
