package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)
var mu sync.Mutex
var count int
func main() {
	http.HandleFunc("/",handler)
	http.HandleFunc("/count",counter)
	log.Fatal(http.ListenAndServe(":8000",nil))
}
func handler(w http.ResponseWriter,r *http.Request) {
	mu.Lock()
	count++
	mu.Unlock()
	fmt.Fprintf(w,"path=%q\n",r.URL.Path)
}
func count(w http.ResponseWriter,r *http.Request) {
	mu.Lock()
	fmt.Fprintf(w,"count:%d\n",count)
	mu.Unlock()
	
}

