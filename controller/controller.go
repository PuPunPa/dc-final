package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"go.nanomsg.org/mangos"

	// register transports
	"go.nanomsg.org/mangos/protocol/pair"
	_ "go.nanomsg.org/mangos/transport/all"
)

type Workload struct {
	ID             int    `json:"ID"`
	Filter         string `json:"Filter"`
	Name           string `json:"Name"`
	Status         string `json:"Status"`
	RunningJobs    int    `json:"RunningJobs"`
	FilteredImages []int  `json:"FilteredImages"`
}

type Image struct {
	ID         int    `json:"ID"`
	WorkloadID int    `json:"WorkloadID"`
	Type       string `json:"Type"`
	Name       string `json:"Name"`
}

type WorkloadsResponse struct {
	Workload []Workload `json:"Workload"`
	Response string     `json:"Response"`
}

var workloads []Workload

var controllerAddress = "tcp://localhost:40899"

var socket mangos.Socket

func die(format string, v ...interface{}) {
	fmt.Fprintln(os.Stderr, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func createWorkload(name string, filter string) bool {
	if filter != "" && filter != "grayscale" && filter != "blur" {
		fmt.Println("False")
		return false
	}
	newWorkload := Workload{len(workloads), filter, name, "Scheduling", 0, make([]int, 0)}
	workloads = append(workloads, newWorkload)
	return true
}

func getWorkload(id int) Workload {
	if id >= len(workloads) {
		return Workload{-1, "", "", "", -1, make([]int, 0)}
	}
	return workloads[id]
}

func getWorkloads() []Workload {
	return workloads
}

func sendMessage(socket mangos.Socket, msg string) {
	if err := socket.Send([]byte(msg)); err != nil {
		die("Controller: There was an error sending the information: %s", err.Error())
	}
}

func receiveMessage(socket mangos.Socket) string {
	var byteArr []byte
	var err error
	if byteArr, err = socket.Recv(); err != nil {
		die("Controller: There was an error receiving the information: %s", err.Error())
	}
	return string(byteArr)
}

func controllerResponses(msg string) {
	response := ""
	var err error
	var rMsg []byte
	split := strings.Split(msg, "_")
	c, _ := strconv.Atoi(split[0])
	split = split[1:]
	switch c {
	//Create Workload, send all workloads
	case 1:
		name := split[0]
		filter := split[1]
		created := createWorkload(name, filter)
		if !created {
			response = ("Controller: There was an error creating the workload")
		}
		respJSON := WorkloadsResponse{workloads, response}
		rMsg, err = json.Marshal(respJSON)
		if err != nil {
			log.Fatal("Controller Error: " + err.Error())
		}

		sendMessage(socket, string(rMsg))
		break

		//send specific workload
	case 2:
		id, _ := strconv.Atoi(split[0])
		workload := getWorkload(id)

		rMsg, err = json.Marshal(workload)
		if err != nil {
			log.Fatal("Controller Error: " + err.Error())
		}

		sendMessage(socket, string(rMsg))
		break
	}
}

func Start() {
	var msg string
	var err error
	if socket, err = pair.NewSocket(); err != nil {
		die("Controller: can't get new pub socket: %s", err)
	}
	if err = socket.Listen(controllerAddress); err != nil {
		die("Controller: can't listen on pub socket: %s", err.Error())
	}
	for {
		msg = receiveMessage(socket)
		controllerResponses(msg)
	}
}
