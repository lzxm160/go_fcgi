 package main
 import (
    _"log"
    "fmt"
    
    "net/http"
)
func pdfHandler (w http.ResponseWriter, r *http.Request) {
  fmt.Fprint(w, "pdf!")
} 
