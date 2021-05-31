package api

import (
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
)

// Logged: Structure to store logged in users info, key is Token, value is User
var online = make(map[string]string)

//Users: stores all users, key Username, value is Password
var users = make(map[string]string)

type Image struct {
	Token, Path, Name string
	size              int
}

// Global variable user. All functions are able to access to it
func Start() {
	router := mux.NewRouter()
	users["username"] = "password"
	//All routes for API
	//Need functions
	router.HandleFunc("/login", login)                         //POST
	router.HandleFunc("/logout", logout)                       //DELETE
	router.HandleFunc("/status", status)                       //GET
	router.HandleFunc("/workloads", workloads)                 //POST
	router.HandleFunc("/workloads/{workload_id}", getWorkload) //GET
	router.HandleFunc("/images", images)                       //POST
	router.HandleFunc("/images/{image_id}", getImage)          //GET
	http.ListenAndServe(":8080", router)
}

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
	if len(token) < 7 {
		w.WriteHeader(http.StatusUnauthorized)
		msg := `
{
	"Please enter a token" ` + token + `smiles
}
`
		w.Write([]byte(msg))
		return
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
	token := r.Header.Get("Authorization")
	if len(token) < 7 {
		w.WriteHeader(http.StatusUnauthorized)
		msg := `
{
	"Please enter a token" ` + token + `
}
`
		w.Write([]byte(msg))
		return
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
		return
	}
	w.WriteHeader(http.StatusOK)
	name, err := os.Hostname()
	if err != nil {
		log.Fatal(err.Error())
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
func workloads(w http.ResponseWriter, r *http.Request) {
	return
}

func getWorkload(w http.ResponseWriter, r *http.Request) {
	return
}

func images(w http.ResponseWriter, r *http.Request) {
	return
}

func getImage(w http.ResponseWriter, r *http.Request) {
	return
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
