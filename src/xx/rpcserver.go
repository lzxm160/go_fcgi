package main

//#go:generate ls -l
import (
	"fmt"
	// "regexp"
	// "os"
	// "bufio"
	// "io"
	//	"io/ioutil"
	"runtime"
	"sync"
	// "testing"
	// "math"
	// "reflect"
	// "unsafe"
	//	"os"
	// "runtime/pprof"
	//	"time"
	// "encoding/json"
	// "bytes"
	// "os/exec"
	// "strings"
	// "syscall"
	// "log"
	"net"
	"net/http"
	"errors"
	"net/rpc"
)

var lock sync.Mutex

func init() {
	// runtime.NumCPU()
	runtime.GOMAXPROCS(runtime.NumCPU())
}

type Args struct{
	A,B int
}
type Quotient struct{
	Quo,Rem int
}
type Arith int
func (t *Arith)Multiply(args *Args,reply *int)error {
	*reply=args.A*args.B
	return nil
}
func (t *Arith)Divide(args *Args,quo *Quotient)error {
	if args.B==0{
		return errors.New("divide by zero")
	}
	quo.Quo=args.A/args.B
	quo.Rem=args.A%args.B
	return nil
}
func main() {
	arith:=new(Arith)
	rpc.Register(arith)
	rpc.HandleHTTP()
	l,e:=net.Listen("tcp",":1234")
	if e!=nil{
		fmt.Println("listen error:",e)
	}
	http.Serve(l,nil)
}
