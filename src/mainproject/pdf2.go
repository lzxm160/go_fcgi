package main

import (
    // _"log"
    "fmt"
    // // "net/http"
    // "unsafe"
    "encoding/json"
    "net/http"
    // "io/ioutil"
    "logger"
    // "strconv"
    // "errors"
    "os"
    "converter"
)
func pdfHandler2 (w http.ResponseWriter, r *http.Request) {

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
	    logger.Info(fmt.Sprintf("Started %s %s for %s:%s response:%s", r.Method, r.URL.Path, addr,body,err_decode.Error()))
	    return
    }
    // logger.Info(t.Src+": "+t.Dst)
    if t.Src==""||t.Dst==""{
    	ret="empty src or dst"
    	fmt.Fprint(w,ret )
    	logger.Info(fmt.Sprintf("Started %s %s for %s:%s response:%s", r.Method, r.URL.Path, addr,body,"empty src or dst"))
	    return
    }
    var err error
    err=convert2(t.Src,t.Dst,ret_convert)
    if err!=nil{
		fmt.Fprint(w,err.Error())
		
    	logger.Info(fmt.Sprintf("Started %s %s for %s:%s response:%s", r.Method, r.URL.Path, addr,body,err.Error()))
    	return
    }else{
    	fmt.Fprint(w,"ok")
    }
    logger.Info(fmt.Sprintf("Started %s %s for %s:%s response:%s", r.Method, r.URL.Path, addr,body,"ok"))
} 
func convert2(src,dst string) error{
	
	source, err := converter.NewConversionSource(src, nil, "pdf")
	if err != nil {
		 
		return err
	}

	return do_convert(*source)
}
func do_convert(source converter.ConversionSource)error {
	// GC if converting temporary file
	if source.IsLocal {
		defer os.Remove(source.URI)
	}

	var conversion converter.Converter
	done := make(chan struct{}, 1)
	got, err := conversion.Convert(source,done)
	if err != nil {
		fmt.Fatalf("convert returned an unexpected error: %+v", err)
		return err
	}
	if want := []byte{}; !reflect.DeepEqual(got, want) {
		fmt.Errorf("expected output of conversion to be %+v, got %+v", want, got)
	}
	return nil
}