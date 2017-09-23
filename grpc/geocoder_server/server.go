package main

import (
	"fmt"
	"log"
	"net"

	"golang.org/x/net/context"

	pb "github.com/reiki4040/go-memo/grpc/address"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct{}

func (s *server) GetAddress(ctx context.Context, in *pb.GeocoderRequest) (*pb.GeocoderReply, error) {
	result := fmt.Sprintf("%f, %f -> %s", in.Lat, in.Lon, "geocoding result dummy")
	return &pb.GeocoderReply{Address: result}, nil
}

func (s *server) GetAddressError(ctx context.Context, in *pb.GeocoderRequest) (*pb.GeocoderReply, error) {
	return nil, fmt.Errorf("dummy geocoding error")
}

func main() {
	l, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterGeocoderServer(s, &server{})

	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(l); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
