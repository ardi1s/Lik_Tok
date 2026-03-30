package observability

import (
	"log"
	"net/http"
	_ "net/http/pprof"
)

type PprofServer struct {
	server *http.Server
}

func NewPprofServer(name string, enabled bool, addr string) (*PprofServer, error) {
	if !enabled {
		return &PprofServer{}, nil
	}

	server := &http.Server{
		Addr: addr,
	}

	go func() {
		log.Printf("Starting %s pprof server on %s", name, addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Pprof server error: %v", err)
		}
	}()

	return &PprofServer{server: server}, nil
}

func (p *PprofServer) Close() error {
	if p.server == nil {
		return nil
	}
	return p.server.Close()
}
