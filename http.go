package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type Server struct {
	// bind address
	ad string
	// connection
	server *http.Server
	// actual KVStore
	kvstore KVStore
}

func NewServer(a string, k KVStore) *Server {
	return &Server{
		ad:      a,
		kvstore: k,
	}
}

func (s *Server) Start() {
	mux := http.NewServeMux()
	mux.HandleFunc("/kv", s.HandleKV)
	mux.HandleFunc("/join", s.HandleJoin)

	s.server = &http.Server{
		Addr:    s.ad,
		Handler: mux,
	}

	go func() {
		err := s.server.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("server closed\n")
		} else if err != nil {
			fmt.Printf("error listening for server: %s\n", err)
		}
	}()
}

func (s *Server) Close() {
	s.server.Close()
}

func (s *Server) HandleKV(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		m := map[string]string{}
		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			w.WriteHeader(http.StatusMethodNotAllowed)
			io.WriteString(w, err.Error())
			return
		}
	case http.MethodPost:
		m := map[string]string{}
		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			w.WriteHeader(http.StatusMethodNotAllowed)
			io.WriteString(w, err.Error())
			return
		}
	case http.MethodDelete:
		m := map[string]string{}
		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			w.WriteHeader(http.StatusMethodNotAllowed)
			io.WriteString(w, err.Error())
			return
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		io.WriteString(w, "Only GET, POST, DELETE allowed.\n")
	}
}

func (s *Server) HandleJoin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		io.WriteString(w, "Join can only be called by a POST operation.\n")
		return
	}
	io.WriteString(w, "JOIN!\n")
}
