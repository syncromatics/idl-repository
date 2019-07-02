package repository

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"

	"github.com/gorilla/mux"
)

type Storage interface {
	ListFolders(path string) ([]string, error)
	File(path string) (io.Reader, error)
	Exists(path string) bool
	MkDir(path string) error
	CreateFile(path string, file io.Reader) error
	ReadFile(path string) (io.ReadCloser, error)
}

type JsonResponse struct {
	StatusCode int
	Model      interface{}
}

type DataResponse struct {
	StatusCode int
	Data       io.ReadCloser
	Error      string
}

type Muxer interface {
	RegisterJson(path string, handler func(HttpContext) (*JsonResponse, error))
	RegisterData(path string, handler func(HttpContext) (*DataResponse, error))
}

type HttpContext struct {
	Args map[string]string
	Body io.Reader
}

type Server struct {
	settings *Settings
	storage  Storage
}

func NewServer(settings *Settings, storage Storage) *Server {
	return &Server{settings, storage}
}

func (s *Server) Run(ctx context.Context) func() error {
	r := mux.NewRouter()
	project := newProjectRouter(s.storage)

	wrap := newRouterWrapper(r)

	project.Register(wrap)

	r.PathPrefix("/").HandlerFunc(s.handle404)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.settings.Port),
		Handler: r,
	}

	cancel := make(chan error)

	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			cancel <- errors.Wrap(err, "failed to serve http")
		}
	}()

	return func() error {
		select {
		case <-ctx.Done():
			// do nothing for now
			return nil
		case msg := <-cancel:
			return msg
		}
	}
}

func (s *Server) handle404(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("url not handled %s\n", r.URL)
	w.WriteHeader(404)
}
