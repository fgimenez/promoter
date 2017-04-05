package service

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"

	"gopkg.in/tylerb/graceful.v1"
)

// Server encapsulates the handler and (graceful) Start/Stop methods
type Server struct {
	srv *graceful.Server
}

// NewServer sets up a negroni server
func NewServer(addr string) *Server {
	n := negroni.Classic()
	mx := mux.NewRouter()

	n.UseHandler(mx)

	s := &Server{
		srv: &graceful.Server{
			Timeout: 2 * time.Second,

			Server: &http.Server{
				Addr:    addr,
				Handler: n,
			},
		},
	}
	return s
}

// Start the server
func (s *Server) Start() error {
	l, err := net.Listen("tcp", s.srv.Addr)
	if err != nil {
		return err
	}

	go s.srv.Serve(l)
	return nil
}

// Stop the server
func (s *Server) Stop() error {
	timeoutTime := 2000 * time.Millisecond
	s.srv.Stop(timeoutTime / 2)

	select {
	case <-s.srv.StopChan():
	case <-time.After(timeoutTime):
		return fmt.Errorf("store failed to stop after %s", timeoutTime)
	}

	return nil
}
