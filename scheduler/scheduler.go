package scheduler

import (
	"context"
	"log"
	"time"
	"C"
	"github.com/HectorJorgeMoralesArch/dc-final/images"
	pb "github.com/CodersSquad/dc-labs/challenges/third-partial/proto"
	"google.golang.org/grpc"
)

//const (
//	address     = "localhost:50051"
//	defaultName = "world"
//)

type Job struct {
	Address string
	RPCName string
}

func schedule(job Job) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(job.Address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	//Blurred Image
	r, err = C.BlurImage(ctx, &pb.Image2Blur{Name: job.RPCName})
	if err != nil {
		log.Fatalf("could not blur the Image: %v", err)
	}
	log.Printf("Scheduler: RPC respose from %s : %s", job.Address, r.GetMessage())
	//GrayScale Image
	r, err = C.GrayScaleImage(ctx, &pb.Image2GrayScale{Name: job.RPCName})
	if err != nil {
		log.Fatalf("could not grayscale the Image: %v", err)
	}
	log.Printf("Scheduler: RPC respose from %s : %s", job.Address, r.GetMessage())
}

func Start(jobs chan Job) error {
	for {
		job := <-jobs
		schedule(job)
	}
	return nil
}
