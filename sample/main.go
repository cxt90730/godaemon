package main

import (
	"fmt"
	daemon "github.com/cxt90730/godaemon"
	"log"
	"net/http"
)

func main() {
	daemon.RunDaemon("mypid.pid", HttpServer)
}

func HttpServer() {
	http.HandleFunc("/", route)
	http.ListenAndServe(":1789", nil)
}
func route(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)
	fmt.Fprint(w, "Hello World\n")
}
