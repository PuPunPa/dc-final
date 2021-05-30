package controller

import (
	"fmt"
	"log"
	"os"
	"time"

	"go.nanomsg.org/mangos"
	"go.nanomsg.org/mangos/protocol/pub"

	// register transports
	_ "go.nanomsg.org/mangos/transport/all"
)

type Workload struct {
	ID             int
	Filter         string
	Name           string
	Status         string
	RunningJobs    int
	FilteredImages []int
}

type Image struct {
	ID         int
	WorkloadID int
	Type       string
	Name       string
}

var workloads []Workload

var controllerAddress = "tcp://localhost:40899"

func die(format string, v ...interface{}) {
	fmt.Fprintln(os.Stderr, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func date() string {
	return time.Now().Format(time.ANSIC)
}

func Start() {
	var sock mangos.Socket
	var err error
	if sock, err = pub.NewSocket(); err != nil {
		die("can't get new pub socket: %s", err)
	}
	if err = sock.Listen(controllerAddress); err != nil {
		die("can't listen on pub socket: %s", err.Error())
	}
	for {
		// Could also use sock.RecvMsg to get header
		d := date()
		log.Printf("Controller: Publishing Date %s\n", d)
		if err = sock.Send([]byte(d)); err != nil {
			die("Failed publishing: %s", err.Error())
		}
		time.Sleep(time.Second * 3)
	}
}

func createWorkload(name string, filter string) (int, bool) {
	newWorkload := Workload{len(workloads), filter, name, "Scheduling", 0, make([]int, 0)}
	workloads = append(workloads, newWorkload)
	return len(workloads), true
}

func getWorkload(id int) Workload {
	return workloads[id]
}

func getWorkloads() []Workload {
	return workloads
}
