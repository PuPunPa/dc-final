package controller

import (
	"fmt"
	"os"
	"time"
	"strings"
	"strconv"
	"io/ioutil"
    "net/http"
    "net/url"
	"github.com/HectorJorgeMoralesArch/dc-final/images"
	"go.nanomsg.org/mangos"
	"go.nanomsg.org/mangos/protocol/pair"

	// register transports
	_ "go.nanomsg.org/mangos/transport/all"
)
// Worker structure
type Worker struct{
	Name string
	Status string
	Usage string
}

// Sockets for survey
var controllerAddress = "tcp://localhost:40899"
var sock mangos.Socket
var err error

// Slice to save registered workers
var Workers []Worker

// Slice to store all workloads
var Workloads []Workload

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

func die(format string, v ...interface{}) {
	fmt.Fprintln(os.Stderr, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func date() string {
	return time.Now().Format(time.ANSIC)
}
func WorkerInfo(name string) (string){
	for _,v := range Workers{
		if v.Name == name{
			return name
		}
	}
	return ""
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
func getActiveWorkloads() int{
	return len(workloads)
}

type Img struct{
	workload_id int `json:"workload_id"`
	image_id int `json:"image_id"`
	Path string `json:"Path"`
	Type string `json:"Type"`
}
var Images []Img
func openNewImage(path string) Img{
	dat, err := ioutil.ReadFile(path)
	if e != nil {
        panic(e)
    }
	newImg :=Img(workload_id: , image_id: len(Images) + 1, Path: path, Type: "orginal")
	Images = append(Images, newImg)
	return newImg
}

func getImage(id int){
	if id>=len(Images){
		w.Write("Image not found")
		return
	}
    // Build fileName from fullPath
    fileURL, err := url.Parse(Images[id].path)
    if err != nil {
        log.Fatal(err)
    }
    path := fileURL.Path
    segments := strings.Split(path, "/")
    fileName = segments[len(segments)-1]
 
    // Create blank file
    file, err := os.Create(fileName)
    if err != nil {
        log.Fatal(err)
    }
    client := http.Client{
        CheckRedirect: func(r *http.Request, via []*http.Request) error {
            r.URL.Opaque = r.URL.Path
            return nil
        },
    }
    // Put content on file
    resp, err := client.Get(fullURLFile)
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()
 
    size, err := io.Copy(file, resp.Body)
 
    defer file.Close()
 
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
	if sock, err = pair.NewSocket(); err != nil {
		die("can't get new pair socket: %s", err)
	}
	if err = sock.Listen(controllerAddress); err != nil {
		die("can't listen on pair socket: %s", err.Error())
	}
	for {
		i := 0
		for {
			var msg []byte
				if msg, err = sock.Recv(); err != nil {
						break
				}
				stats := strings.Split(string(msg), " ")
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
