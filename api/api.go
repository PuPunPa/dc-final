package api

import (
	"encoding/json"
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	//If it says that Go could not found the mrepository just type
	//go get github.com/gorilla/mux
	//to fix it
	"github.com/gorilla/mux"
	"go.nanomsg.org/mangos"
	"go.nanomsg.org/mangos/protocol/pair"
)

type Workload struct {
	ID             int    `json:"ID"`
	Filter         string `json:"Filter"`
	Name           string `json:"Name"`
	Status         string `json:"Status"`
	RunningJobs    int    `json:"RunningJobs"`
	FilteredImages []int  `json:"FilteredImages"`
}

type WorkloadsResponse struct {
	Workload []Workload
	Response string
}

type Response struct {
	Workload Workload
	Response string
}

// Logged: Structure to store logged in users info, key is Token, value is User
var online = make(map[string]string)

//Users: stores all users, key Username, value is Password
var users = make(map[string]string)

type Image struct {
	Token, Path, Name string
	size              int
}

var socket mangos.Socket

var controllerAddress = "tcp://localhost:40899"
var apiAddress = "tcp://localhost:40900"

func login(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-type", "application/json")
	usr, pwd, ok := r.BasicAuth()
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		msg := `
{
	"Username and Password are required"
}
`
		w.Write([]byte(msg))
		return
	}
	if !checkPass(usr, pwd) {
		w.WriteHeader(http.StatusUnauthorized)
		msg := `
{
	"Invalid Username and/or Password"
}
`
		w.Write([]byte(msg))
		return
	}
	w.WriteHeader(http.StatusOK)
	token := GenerateRandomString(10)
	msg := `
{
	"message": "Hi ` + usr + ` welcome to the DPIP System"
	"token": ` + token + `"
}
`
	w.Write([]byte(msg))
	online[token] = usr
	return
}

func checkPass(usr string, pwd string) bool {
	if p, ok := users[usr]; ok {
		return p == pwd
	}
	return false
}

func logout(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if !checkToken(w, r) {
		return
	}
	w.WriteHeader(http.StatusOK)
	msg := `
{
	"logout_message": "Bye ` + online[token] + `, your token has been revoked"
}
`
	w.Write([]byte(msg))
	delete(online, token)
	return
}

func loggedIn(token string) bool {
	if _, ok := online[token]; ok {
		return true
	}
	return false
}

func status(w http.ResponseWriter, r *http.Request) {
	if !checkToken(w, r) {
		return
	}
	w.WriteHeader(http.StatusOK)
	name, err := os.Hostname()
	if err != nil {
		log.Fatal("API Error: " + err.Error())
	}
	msg := `
{
	"system_name": ` + name + `"
	"server_time": "` + time.Now().Format("2006-01-02 15:04:05") + `"
	"active_workloads": "
}
`
	w.Write([]byte(msg))
	return
}

func images(w http.ResponseWriter, r *http.Request) {
	if !checkToken(w, r) {
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}

func workloads(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-type", "application/json")
	if !checkToken(w, r) {
		return
	}
	w.WriteHeader(http.StatusOK)
	name := r.Header.Get("WorkloadName")
	filter := r.Header.Get("Filter")
	sendMessage(socket, "1_"+name+"_"+filter)
	rMsg := receiveMessage(socket)
	workloadResp := new(WorkloadsResponse)
	err := json.Unmarshal([]byte(rMsg), &workloadResp)
	if err != nil {
		log.Fatal("API Error: " + err.Error())
	}
	workloads := workloadResp.Workload
	respResp := workloadResp.Response
	if respResp != "" {
		w.Write([]byte(respResp))
		return
	}
	for _, wl := range workloads {
		msg, err := json.MarshalIndent(wl, " ", "    ")
		if err != nil {
			log.Fatal("API error: " + err.Error())
		}
		w.Write([]byte(msg))
		w.Write([]byte("\n"))
	}
	return
}

func getWorkload(w http.ResponseWriter, r *http.Request) {
	if !checkToken(w, r) {
		return
	}
	w.WriteHeader(http.StatusOK)
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	id := vars["workload_id"]
	sendMessage(socket, "2_"+id)
	rMsg := receiveMessage(socket)
	workloadResp := new(Workload)
	err := json.Unmarshal([]byte(rMsg), &workloadResp)
	if err != nil {
		log.Fatal("API Error: " + err.Error())
	}
	msg, err := json.MarshalIndent(workloadResp, " ", "    ")
	if err != nil {
		log.Fatal("API error: " + err.Error())
	}
	w.Write([]byte(msg))
	w.Write([]byte("\n"))
	return
}

func getImage(w http.ResponseWriter, r *http.Request) {
	if !checkToken(w, r) {
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}

func checkToken(w http.ResponseWriter, r *http.Request) bool {
	token := r.Header.Get("Authorization")
	if len(token) < 7 {
		w.WriteHeader(http.StatusUnauthorized)
		msg := `
{
	"Please enter a token" ` + token + `
}
`
		w.Write([]byte(msg))
		return false
	}
	token = token[7:]
	if !loggedIn(token) {
		w.WriteHeader(http.StatusUnauthorized)
		msg := `
{
	"Please enter a valid token"
}
`
		w.Write([]byte(msg))
		return false
	}
	return true
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// GenerateRandomString generate a string of random characters of given length
func GenerateRandomString(n int) string {
	sb := strings.Builder{}
	sb.Grow(n)
	for i := 0; i < n; i++ {
		idx := rand.Int63() % int64(len(letterBytes))
		sb.WriteByte(letterBytes[idx])
	}
	return sb.String()
}

func sendMessage(socket mangos.Socket, msg string) {
	if err := socket.Send([]byte(msg)); err != nil {
		die("API: There was an error sending the information: %s", err.Error())
	}
}

func receiveMessage(socket mangos.Socket) string {
	var byteArr []byte
	var err error
	if byteArr, err = socket.Recv(); err != nil {
		die("API: There was an error receiving the information: %s", err.Error())
	}
	return string(byteArr)
}

func die(format string, v ...interface{}) {
	fmt.Fprintln(os.Stderr, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func Start() {
	var err error
	if socket, err = pair.NewSocket(); err != nil {
		die("API: can't get new pub socket: %s", err)
	}
	if err = socket.Dial(controllerAddress); err != nil {
		die("API: can't listen on pub socket: %s", err.Error())
	}
	router := mux.NewRouter()
	users["username"] = "password"

	router.HandleFunc("/login", login)                         //POST
	router.HandleFunc("/logout", logout)                       //DELETE
	router.HandleFunc("/status", status)                       //GET
	router.HandleFunc("/workloads", workloads)                 //POST
	router.HandleFunc("/workloads/{workload_id}", getWorkload) //GET
	router.HandleFunc("/images", images)                       //POST
	router.HandleFunc("/images/{image_id}", getImage)          //GET
	http.ListenAndServe(":8080", router)
}
