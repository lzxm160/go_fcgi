package main
//#cgo CFLAGS: -I../wkhtmltox/include 
//#cgo LDFLAGS: -L/usr/lib -lwkhtmltox -Wall -ansi -pedantic -ggdb
//#include <stdbool.h>
//#include <stdio.h>
//#include <string.h>
//#include <stdlib.h>
//#include "../wkhtmltox/include/pdf.h"
//extern void finished_cb(void*, const int);
//extern void progress_changed_cb(void*, const int);
//extern void error_cb(void*, char *msg);
//extern void warning_cb(void*, char *msg);
//extern void phase_changed_cb(void*);
//static void setup_callbacks(wkhtmltopdf_converter * c) {
//  wkhtmltopdf_set_finished_callback(c, (wkhtmltopdf_int_callback)finished_cb);
//  wkhtmltopdf_set_progress_changed_callback(c, (wkhtmltopdf_int_callback)progress_changed_cb);
//  wkhtmltopdf_set_error_callback(c, (wkhtmltopdf_str_callback)error_cb);
//  wkhtmltopdf_set_warning_callback(c, (wkhtmltopdf_str_callback)warning_cb);
//  wkhtmltopdf_set_phase_changed_callback(c, (wkhtmltopdf_void_callback)phase_changed_cb);
//}
import "C"
import (
    _"log"
    "fmt"
    // "net/http"
    "unsafe"
    "encoding/json"
    "net/http"
    "io/ioutil"
    "logger"
    "strconv"
    "errors"
)
type GlobalSettings struct {
	s *C.wkhtmltopdf_global_settings
}

type ObjectSettings struct {
	s *C.wkhtmltopdf_object_settings
}

type Converter struct {
	c               *C.wkhtmltopdf_converter
	Finished        func(*Converter, int)
	ProgressChanged func(*Converter, int)
	Error           func(*Converter, string)
	Warning         func(*Converter, string)
	Phase           func(*Converter)
}

var converter_map map[unsafe.Pointer]*Converter

func NewGolbalSettings() *GlobalSettings {
	return &GlobalSettings{s: C.wkhtmltopdf_create_global_settings()}
}

func (self *GlobalSettings) Set(name, value string) {
	c_name := C.CString(name)
	c_value := C.CString(value)
	defer C.free(unsafe.Pointer(c_name))
	defer C.free(unsafe.Pointer(c_value))
	C.wkhtmltopdf_set_global_setting(self.s, c_name, c_value)
}
func (self *GlobalSettings) SetBool(name string, value bool) {
	c_name := C.CString(name)
	c_value := C.bool(value)
	defer C.free(unsafe.Pointer(c_name))
	defer C.free(unsafe.Pointer(c_value))
	C.wkhtmltopdf_set_global_setting(self.s, c_name, c_value)
}
func NewObjectSettings() *ObjectSettings {
	return &ObjectSettings{s: C.wkhtmltopdf_create_object_settings()}
}

func (self *ObjectSettings) Set(name, value string) {
	c_name := C.CString(name)
	c_value := C.CString(value)
	defer C.free(unsafe.Pointer(c_name))
	defer C.free(unsafe.Pointer(c_value))
	C.wkhtmltopdf_set_object_setting(self.s, c_name, c_value)
}
func (self *ObjectSettings) SetBool(name string, value bool) {
	c_name := C.CString(name)
	c_value := C.bool(value)
	defer C.free(unsafe.Pointer(c_name))
	defer C.free(unsafe.Pointer(c_value))
	C.wkhtmltopdf_set_object_setting(self.s, c_name, c_value)
}
func (self *GlobalSettings) NewConverter() *Converter {
	c := &Converter{c: C.wkhtmltopdf_create_converter(self.s)}
	C.setup_callbacks(c.c)

	return c
}
func (self *Converter) Convert() error {

	// To route callbacks right, we need to save a reference
	// to the converter object, base on the pointer.
	converter_map[unsafe.Pointer(self.c)] = self
	status := C.wkhtmltopdf_convert(self.c)
	delete(converter_map, unsafe.Pointer(self.c))
	if status != C.int(1) {
		return errors.New("Convert failed")
	}
	// fmt.Printf("status: %d\n", status)
	return nil
}
func convert(src,dst string) error {
	converter_map = make(map[unsafe.Pointer]*Converter)
	C.wkhtmltopdf_init(C.false)
	gs := NewGolbalSettings()
	gs.Set("outputFormat", "pdf")
	gs.Set("out", dst)
	gs.Set("orientation", "Portrait")
	gs.Set("colorMode", "Color")
	gs.Set("size.paperSize", "A4")
	//gs.Set("load.cookieJar", "myjar.jar")
	// object settings: http://www.cs.au.dk/~jakobt/libwkhtmltox_0.10.0_doc/pagesettings.html#pagePdfObject
	os := NewObjectSettings()
	os.Set("page", src)
	// os.Set("load.debugJavascript", "false")
	//os.Set("load.jsdelay", "1000") // wait max 1s
	os.Set("web.enableJavascript", "true")
	os.Set("web.enablePlugins", "true")
	os.Set("web.loadImages", "true")
	os.Set("web.background", "true")
	os.Set("web.defaultEncoding", "utf-8")
	// os.Set("web.userStyleSheet", "utf-8")
	// os.Set("load.blockLocalFileAccess","false")
	os.Set("load.blockLocalFileAccess",false) 
	os.Set("load.loadErrorHandling","skip")
	c := gs.NewConverter()
	c.Add(os)
	//c.AddHtml(os, "<html><body><h3>HELLO</h3><p>World</p></body></html>")

	c.ProgressChanged = func(c *Converter, b int) {
		// fmt.Printf("Progress: %d\n", b)
	}
	c.Error = func(c *Converter, msg string) {
		// fmt.Printf("error: %s\n", msg)
		logger.Error("error: "+msg)
            
	}
	c.Warning = func(c *Converter, msg string) {
		// fmt.Printf("warning: %s\n", msg)
		logger.Warn("warning: " + msg)
	}
	c.Phase = func(c *Converter) {
		// fmt.Printf("Phase\n")
	}
	c.Finished = func(c *Converter, s int) {
		// fmt.Printf("Finished: %d\n", s)
		logger.Info("pdf Finished:" + strconv.Itoa(s))
	}
	err:=c.Convert()
	// temp:=c.ErrorCode()
	// logger.Info("Got error code: " + strconv.Itoa(temp))
	// fmt.Printf("Got error code: %d\n", temp)

	c.Destroy()
	C.wkhtmltopdf_deinit()	
	if err!=nil{
		return err
	}
	
	return nil
}
type src_dst struct{
	Src string
	Dst string
}
func pdfHandler (w http.ResponseWriter, r *http.Request) {
  	addr := r.Header.Get("X-Real-IP")
	if addr == "" {
		addr = r.Header.Get("X-Forwarded-For")
		if addr == "" {
			addr = r.RemoteAddr
		}
	}
  
	body, _:= ioutil.ReadAll(r.Body)
	// logger.Info(fmt.Sprintf("%s",body))
	defer r.Body.Close()
    var ret string
    var t src_dst  
    err_decode := json.Unmarshal(body, &t)
    if err_decode!=nil{
    	ret=`decode failed`
	    fmt.Fprint(w,ret )
	    return
    }
    // logger.Info(t.Src+": "+t.Dst)
    if t.Src==""||t.Dst==""{
    	ret="empty src or dst"
    	fmt.Fprint(w,ret )
	    return
    }
    var err error
    err=convert(t.Src,t.Dst)
    if err!=nil{
		fmt.Fprint(w,err.Error())
		log_str:=fmt.Sprintf("Started %s %s for %s:%s response:%s", r.Method, r.URL.Path, addr,body,err.Error())
    	logger.Info(log_str)
    	return
    }else{
    	fmt.Fprint(w,"ok")
    }
	log_str:=fmt.Sprintf("Started %s %s for %s:%s response:%s", r.Method, r.URL.Path, addr,body,"ok")
    logger.Info(log_str)
} 









///////////////////////////////////////////////////
func (self *Converter) Add(settings *ObjectSettings) {
	C.wkhtmltopdf_add_object(self.c, settings.s, nil)
}

func (self *Converter) AddHtml(settings *ObjectSettings, data string) {
	c_data := C.CString(data)
	defer C.free(unsafe.Pointer(c_data))
	C.wkhtmltopdf_add_object(self.c, settings.s, c_data)
}

func (self *Converter) ErrorCode() int {
	return int(C.wkhtmltopdf_http_error_code(self.c))
}

func (self *Converter) Destroy() {
	C.wkhtmltopdf_destroy_converter(self.c)
}

//export finished_cb
func finished_cb(c unsafe.Pointer, s C.int) {
	conv := converter_map[c]
	if conv.Finished != nil {
		conv.Finished(conv, int(s))
	}
}

//export progress_changed_cb
func progress_changed_cb(c unsafe.Pointer, p C.int) {
	conv := converter_map[c]
	if conv.ProgressChanged != nil {
		conv.ProgressChanged(conv, int(p))
	}
}

//export error_cb
func error_cb(c unsafe.Pointer, msg *C.char) {
	conv := converter_map[c]
	if conv.Error != nil {
		conv.Error(conv, C.GoString(msg))
	}
}

//export warning_cb
func warning_cb(c unsafe.Pointer, msg *C.char) {
	conv := converter_map[c]
	if conv.Warning != nil {
		conv.Warning(conv, C.GoString(msg))
	}
}

//export phase_changed_cb
func phase_changed_cb(c unsafe.Pointer) {
	conv := converter_map[c]
	if conv.Phase != nil {
		conv.Phase(conv)
	}
}