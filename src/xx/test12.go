package main
//#go:generate ls -l
import (
	"fmt"
	// "regexp"
	"os"
	// "bufio"
	"runtime"
	"io"
	"io/ioutil"
	"sync"
	// "testing"
	// "math"
	// "reflect"
	// "unsafe"
	// "os"
	// "runtime/pprof"
	"time"
	// "encoding/json"
	// "bytes"
	"os/exec"
	"strings"
)

var lock sync.Mutex

func init() {
	// runtime.NumCPU()
	runtime.GOMAXPROCS(runtime.NumCPU())
}
func bind(f func(in io.Reader,out io.Writer,p []string),params []string)func(in io.Reader,out io.Writer) {
	return func(in io.Reader,out io.Writer) {
		f(in,out,params)
	}
}
func pipe(app1 func(in io.Reader,out io.Writer),app2 func(in io.Reader,out io.Writer))func(in io.Reader,out io.Writer) {
	return func(in io.Reader,out io.Writer) {
		r,w:=io.Pipe()
		defer w.Close()
		go func() {
			defer r.Close()
			app2(r,out)
		}()
		app1(in,w)
	}
}
func main() {
	func1:=func(in io.Reader,out io.Writer,params []string) {
		var command string
		if b, err := ioutil.ReadAll(in); err == nil {
		    command=string(b)
		}
		cmd := exec.Command(command, params[0])
		cmd.Stdout=out
		err := cmd.Start()
		if err != nil {
			fmt.Println("cmd error!")
		}
	}
	params1:=[]string{"test"}
	bind1:=bind(func1,params1)

	func2:=func(in io.Reader,out io.Writer,params []string) {
		// var content string
		// if b, err := ioutil.ReadAll(in); err == nil {
		//     content=string(b)
		// }else{
		// 	fmt.Println("content:",err.Error())
		// }
		// fmt.Println("content:",content)
		cmd := exec.Command("grep", params[0],in)
		cmd.Stdout =out
		err := cmd.Start()
		if err != nil {
			fmt.Println("cmd error!")
		}
	}
	params2:=[]string{"select"}
	bind2:=bind(func2,params2)
	// fp(strings.NewReader("cat"),os.Stdout)
	fmt.Println("--------------------")
	pp:=pipe(bind1,bind2)
	pp(strings.NewReader("cat"),os.Stdout)
	//cat note|grep select
	time.Sleep(2*time.Second)
	fmt.Println("done!")
}


