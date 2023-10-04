package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"time"

	pb "github.com/twonull/grpc-inspect/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewInspectClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	r, err := c.SendInspect(ctx, &pb.InspectRequest{
		// Example Inspect Request -- Five-SeveN | Berries and Cherries (Factory New)
		Fields: []uint64{76561198181851863, 20029930164, 10234124678058918143, 76561198181851863},
	})
	if err != nil {
		log.Fatalf("could not send req: %v", err)
	}
	// Raw ItemInfo data returned from Steam Protobufs. Fields with no value ignored.
	log.Println("Item Info:", r.GetItemInfo())
	// Wear value converted from int to Float32
	log.Println("Float Value:", getFv(r.GetItemInfo().GetPaintWear()))
}

func getFv(f uint32) string {
	fv := math.Float32frombits(uint32(f))
	return fmt.Sprintf("%.25f", fv)
}
