package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"image"
	"image/color"
	"image/gif"
	"io"
	"math"
	"math/rand"
)
var palette=[]color.Color{color.White,color.Black}
const(
	whiteIndex=0
	blackIndex=1
)
var mu sync.Mutex
var count int
func main() {
	type Celsius float64
	type Fahrenheit float64
	var c Celsius
	var f Fahrenheit
	fmt.Println(c==0)
	fmt.Println(f>=0)
	fmt.Println(c==f)
	fmt.Println(c==Celsius(f))
	http.HandleFunc("/h",handler)
	http.HandleFunc("/count",counter)
	http.HandleFunc("/l",l)
	log.Fatal(http.ListenAndServe(":8000",nil))
}
func handler(w http.ResponseWriter,r *http.Request) {
	mu.Lock()
	count++
	mu.Unlock()
	fmt.Fprintf(w,"path=%q\n",r.URL.Path)
}
func counter(w http.ResponseWriter,r *http.Request) {
	mu.Lock()
	fmt.Fprintf(w,"count:%d\n",count)
	mu.Unlock()
}
func l(w http.ResponseWriter,r *http.Request) {
	lissajous(w)
}
func lissajous(out io.Writer) {
	const(
		cycles=5
		res=0.001
		size=100
		nframes=64
		delay=10
	)
	freq:=rand.Float64()*3.0
	anim:=gif.GIF{LoopCount: nframes}
	phase:=0.0
	for i:=0;i<nframes;i++{
		rect:=image.Rect(0,0,2*size+1,2*size+1)
		// image.NewPalletted(rect,palette)
		img:=image.NewPaletted(rect, palette)
		for t:=0.0;t<cycles*2*math.Pi;t+=res{
			x:=math.Sin(t)
			y:=math.Sin(t*freq+phase)
			img.SetColorIndex(size+int(x*size+0.5),size+int(y*size+0.5),blackIndex)
		}
		phase+=0.1
		anim.Delay=append(anim.Delay,delay)
		anim.Image=append(anim.Image,img)
	}
	gif.EncodeAll(out,&anim)
}