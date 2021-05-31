package controller

import (
	"fmt"
	"log"
	"os"
	"time"
	"strconv"
	"go.nanomsg.org/mangos"
	"go.nanomsg.org/mangos/protocol/pub"

	// register transports
	_ "go.nanomsg.org/mangos/transport/all"
)



// Worker structure
type Worker struct{
	Name string
	Status string
	Usage string
}
// Workload structure
type Workload struct{
	Name string
	Id int
}

// Sockets for survey
var controllerAddress = "tcp://localhost:40899"
var sock mangos.Socket
var err error

// Slice to save registered workers
var Workers []Worker

// Slice to store all workloads
var Workloads []Workload


func die(format string, v ...interface{}) {
	fmt.Fprintln(os.Stderr, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func date() string {
	return time.Now().Format(time.ANSIC)
}
func WorkerInfo(name string) (string)
{
	for _,v := range Workers{
		if v.Name == name{
			return 
		}
	}
}

//Function to add workload
func WorkloadId(name string) int{
	exists := false
	for i,v := range Workloads{
			if v.Name == name{
				exists = true
					allWorkloads[i].jobId++
					return allWorkloads[i].jobId
			}
	}
	if !exists{
			newWorkload := workload{Name: name, Id: 1}
			Workloads = append(Workloads, newWorkload)
			return 1
	}
	return 0
}

// Check worker info
func WorkerStatus(name string) (string, string, string){
	for _,v := range Workers{
			if v.Name == name{
					return v.Name, v.Status, v.Usage
			}
	}
	return "", "", ""
}


func Start() {
	if sock, err = pub.NewSocket(); err != nil {
		die("can't get new pub socket: %s", err)
	}
	if err = sock.Listen(controllerAddress); err != nil {
		die("can't listen on pub socket: %s", err.Error())
	}
	for {
		fmt.Println("Checking workers")
		d := date()
		if err = sock.Send([]byte(d)); err != nil {
			die("Failed: %s", err.Error())
		}
		i := 0
		for {
				if msg, err = sock.Recv(); err != nil {
						break
				}
				stats := strings.Split(string(msg), ",")
				name := "Worker " + strconv.Itoa(i)
				newWorker := Worker{Name: name, Status: stats[0], Usage: stats[1]}
				alive := false
				for _,v := range Workers{
						if v.Name == name{
							alive = true
						}
				}
				if !alive{
						Workers = append(Workers, newWorker)
				}
				fmt.Printf("%s is alive\n", name)
				i++
		}
		time.Sleep(time.Second * 3)
	}
}
