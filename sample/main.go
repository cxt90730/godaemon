package main

import (
	"net/http"
	"log"
	"fmt"
)

func main() {

}

func HttpServer() {
	http.HandleFunc("/", route)
	http.ListenAndServe(":1789", nil)
}
func route(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)
	fmt.Fprint(w, "Hello World\n")
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

