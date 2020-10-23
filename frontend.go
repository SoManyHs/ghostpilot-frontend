package frontend

import (
	"flag"
	"frontend/server"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func Run() error {
	addr := flag.String("addr", ":8080", "port to listen on")
	flag.Parse()

	s := http.Server{
		Addr: *addr,
		Handler: &server.Server{
			Router: mux.NewRouter(),
		},
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}
	log.Printf("INFO: vote: listen on port %s\n", *addr)
	return s.ListenAndServe()
}