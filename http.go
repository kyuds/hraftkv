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
		server:  nil,
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
		key := r.URL.Query().Get("key")
		if len(key) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, "Operation needs a key query.\n")
			return
		}
		ret, err := s.kvstore.Get(key)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, err.Error())
		} else {
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, fmt.Sprintf("%s\n", ret))
		}

	case http.MethodPost:
		m := map[string]string{}
		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, err.Error())
			return
		}
		if len(m) == 0 {
			w.WriteHeader(http.StatusOK)
			return
		}
		err := s.kvstore.Put(m)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, err.Error())
			return
		}
		w.WriteHeader(http.StatusOK)

	case http.MethodDelete:
		key := r.URL.Query().Get("key")
		if len(key) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, "Operation needs a key query.\n")
			return
		}
		err := s.kvstore.Delete(key)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, err.Error())
		} else {
			w.WriteHeader(http.StatusOK)
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

	m := map[string]string{}
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, err.Error())
		return
	}

	addr, ok := m["addr"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "Join must provide address")
		return
	}
	id, ok := m["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "Join must provide node id")
		return
	}

	w.WriteHeader(http.StatusOK)
	io.WriteString(w, fmt.Sprintf("Joining <%s> on %s\n", id, addr))
}

func (s *Server) Addr() string {
	return s.ad
}
