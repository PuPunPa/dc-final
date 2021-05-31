package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"
	"github.com/HectorJorgeMoralesArch/dc-final/api"
	"github.com/HectorJorgeMoralesArch/dc-final/controller"
	"github.com/HectorJorgeMoralesArch/dc-final/scheduler"
	"github.com/HectorJorgeMoralesArch/dc-final/images"
	"github.com/PuPunPa/dc-final/api"
	"github.com/PuPunPa/dc-final/controller"
	"github.com/PuPunPa/dc-final/scheduler"
	"github.com/PuPunPa/dc-final/images"
)

func main() {
	log.Println("Welcome to the Distributed and Parallel Image Processing System")

	// Start Controller
	go controller.Start()

	// Start Scheduler
	jobs := make(chan scheduler.Job)
	go scheduler.Start(jobs)
	// Send sample jobs
	sampleJob := scheduler.Job{Address: "localhost:50051", RPCName: "hello"}

	
	// API
	// Here's where your API setup will be
	go api.Start()

	for {
		sampleJob.RPCName = fmt.Sprintf("hello-%v", rand.Intn(10000))
		jobs <- sampleJob
		time.Sleep(time.Second * 5)
	}
}
