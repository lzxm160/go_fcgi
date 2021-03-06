package main

import (
	"fmt"
	"regexp"
	"os"
	"bufio"
	"runtime"
	"io"
)
type safeMap chan commandData
type commandData struct{
	action commandAction
	key string
	value interface{}
	result chan<- interface{}
	data chan<- map[string]interface{}
	updater UpdateFunc
}
type commandAction int
const (
	remove commandAction=iota
	end 
	find
	insert
	length
	update
)
type findResult struct{
	value interface{}
	found bool
}
type SafeMap interface{
	Insert(string,interface{})
	Delete(string)
	Find(string)(interface{},bool)
	Len()int
	Update(string,UpdateFunc)
	Close()map[string]interface{}
}
type UpdateFunc func(interface{},bool)interface{}
func New()(SafeMap) {
	sm:=make(safeMap)
	go sm.run()
	return sm
}
func (sm safeMap)run() {
	store:=make(map[string]interface{})
	for command:=range sm{
		switch command.action{
		case insert:
			store[command.key]=command.value
		case remove:
			delete(store,command.key)
		case find:
			value,found:=store[command.key]
			command.result<-findResult{value,found}
		case length:
			command.result<-len(store)
		case update:
			value,found:=store[command.key]
			store[command.key]=command.updater(value,found)
		case end:
			close(sm)
			command.data<-store
		}
	}
}
func (sm safeMap)Insert(key string,value interface{}) {
	(sm)<-commandData{action:insert,key:key,value:value}
}
func (sm safeMap)Delete(key string) {
	(sm)<-commandData{action:remove,key:key}
}

func (sm safeMap)Find(key string)(value interface{},found bool) {
	reply:=make(chan interface{})
	(sm)<-commandData{action:find,key:key,result:reply}
	result:=(<-reply).(findResult)
	return result.value,result.found
}
func (sm safeMap)Len()int {
	reply:=make(chan interface{})
	(sm)<-commandData{action:length,result:reply}
	return (<-reply).(int)
}
func (sm safeMap)Update(key string,updater UpdateFunc){
	(sm)<-commandData{action:update,key:key,updater:updater}
}
func (sm safeMap)Close()map[string]interface{} {
	reply:=make(chan map[string]interface{})
	(sm)<-commandData{action:end,data:reply}
	return <-reply
}

func main() {
	// test :=New()
	// test.Insert("1",2)
	// if data,found:=test.Find("1");found{
	// 	fmt.Println(data)
	// }
	filename:="/root/redisRenesola-cluster-debug/logs/cache_20170227_00209.log"
	var workers=runtime.NumCPU()
	fmt.Println(workers)
	runtime.GOMAXPROCS(workers)
	lines:=make(chan string,32)
	done:=make(chan struct{},8)
	pageMap:=New()
	// for i:=0;i<20;i++{
	// 	pageMap.Insert(fmt.Sprintf("%d",i),i)
	// }
	// show(pageMap)
	go readlines(filename,lines)
	processLines(done,pageMap,lines)
	waitUntil(done)
	show(pageMap)
}
func readlines(filename string,lines chan<- string) {
	file,err:=os.Open(filename)
	if err!=nil{
		fmt.Println(err)
	}
	defer file.Close()
	reader:=bufio.NewReader(file)
	for{
		line,err:=reader.ReadString('\n')
		if line!=""{
			lines<-line
		}
		if err!=nil{
			if err!=io.EOF{
				fmt.Println(err)
			}
			break
		}
	}
	close(lines)
}
func processLines(done chan<- struct{},pageMap SafeMap,lines <-chan string){
	reg:=regexp.MustCompile("GET /flowNo/.?")
	incrementer:=func(value interface{},found bool)interface{} {
		if found{
			return value.(int)+1
		}
		return 1
	}
	for i:=0;i<8;i++{
		go func(){
			for line:=range lines{
				if matches:=reg.FindStringSubmatch(line);matches!=nil{
					pageMap.Update(matches[0],incrementer)
					// fmt.Println(matches)
				}
			}
			done<-struct{}{}
		}()

	}
}
func waitUntil(done <-chan struct{}) {
	for i:=0;i<8;i++{
		<-done
		fmt.Printf("160:%d\n",i)
	}
}
func show(pm SafeMap) {
	pages:=pm.Close()
	for page,count:=range pages{
		fmt.Printf("%8d %s\n",count,page)
	}
}