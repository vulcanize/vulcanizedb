package main

import (
	"log"
	"net/http"

	"github.com/fjl/memsize/memsizeui"
)

func main() {
	byteslice := make([]byte, 200)
	intslice := make([]int, 100)

	h := new(memsizeui.Handler)
	s := &http.Server{Addr: "127.0.0.1:8080", Handler: h}
	h.Add("byteslice", &byteslice)
	h.Add("intslice", &intslice)
	h.Add("server", s)
	log.Println("listening on http://127.0.0.1:8080")
	s.ListenAndServe()
}
