package main

import (
	"log"

	"golang.org/x/net/context"

	pb "github.com/reiki4040/go-memo/grpc/address"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:3000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewGeocoderClient(conn)

	r, err := c.GetAddress(context.Background(), &pb.GeocoderRequest{Lat: 35, Lon: 139})
	if err != nil {
		errc := grpc.Code(err)
		log.Printf("error code: %d", errc)
		log.Fatalf("could not geocode: %v", err)
	}
	log.Printf("Geocoder result: %s", r.Address)

	r, err = c.GetAddressError(context.Background(), &pb.GeocoderRequest{Lat: 35, Lon: 139})
	if err != nil {
		errc := grpc.Code(err)
		log.Printf("error code: %d", errc)
		log.Fatalf("could not geocode: %v", err)
	}
	log.Printf("Geocoder result: %s", r.Address)
}
