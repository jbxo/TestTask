package shared

import (
	"errors"
	"net/http"
)

var (
	defaultParameters = &ServerParameters{
		AddressToListen: ":8080",
	}
)

type Endpoint struct {
	Adapters []Adapter
	Handler  http.Handler
	Path     string
}

type ServerParameters struct {
	AddressToListen string
}

type Server struct {
	Parameters *ServerParameters
	hServer    *http.Server
	mux        *http.ServeMux
	endpoints  []Endpoint
}

func NewServer(parameters *ServerParameters, endpoints ...Endpoint) *Server {
	if parameters == nil {
		parameters = defaultParameters
	}

	return &Server{
		Parameters: parameters,
		endpoints:  endpoints,
	}
}

func (s *Server) AddEndpoints(es ...Endpoint) {
	for _, e := range es {
		if e.Handler == nil {
			continue
		}
		s.endpoints = append(s.endpoints, e)
	}
}

func buildMuxFromEndpoints(es []Endpoint) *http.ServeMux {
	mux := http.NewServeMux()

	for _, e := range es {
		mux.Handle(e.Path, adaptersHandlerExecutor(e.Handler, e.Adapters))
	}

	return mux
}

func (s *Server) Start() error {
	if len(s.endpoints) == 0 {
		errors.New("don't know how to start server without endpoints")
	}
	s.mux = buildMuxFromEndpoints(s.endpoints)

	s.hServer = &http.Server{
		Addr:    s.Parameters.AddressToListen,
		Handler: s.mux,
	}

	return s.hServer.ListenAndServe()
}
