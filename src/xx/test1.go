package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func main() {
	// lissajous(os.Stdout)
	test_log()
}
func test_log() {
	start:=time.Now()
	ch:=make(chan string)
	for _,url:=range os.Args[1:]{
		go fetch(url,ch)
	}
	for range os.Args[1:]{
		fmt.Println(<-ch)
	}
	fmt.Printf("%.2fs elapsed\n",time.Since(start).Seconds())
}
func fetch(url string,ch chan<-string) {
	start:=time.Now()
	resp,err:=http.Get(url)
	if err!=nil{
		secs:=time.Since(start).Seconds()
		fmt.Sprintf("%.2f %s:%s",secs,url)
		// ch<-fmt.Sprint(temp+":"+err)
		ch<-fmt.Sprintf("%.2f %s:%s",secs,url,err)
		return
	}
	nbytes,err:=io.Copy(ioutil.Discard,resp.Body)
	resp.Body.Close()
	if err!=nil{
		ch<-fmt.Sprintf("reading err:%s %v",url,err)
		return
	}
	secs:=time.Since(start).Seconds()
	ch<-fmt.Sprintf("%.2f %7d %s",secs,nbytes,url)
}
