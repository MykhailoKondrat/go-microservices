package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

type Hello struct {
	l *log.Logger
}
type Goodbye struct {
	l *log.Logger
}

func NewHello(l *log.Logger) *Hello {
	return &Hello{l}
}
func NewGoodBuy(l *log.Logger) *Goodbye {
	return &Goodbye{l}
}

func (g *Goodbye) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Byeee"))
	if err != nil {
		panic(err)
	}
}

func (h *Hello) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.l.Println("hello world")
	data, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, "Hello %s", data)
}
