package internalhttp

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Server struct {
	Router *mux.Router
	Application
	l *zap.Logger
}

type Application interface{}

func NewServer(app Application, log *zap.Logger) *Server {
	// nolint:exhaustivestruct
	return &Server{
		Application: app,
		l:           log,
	}
}

func (s *Server) Start(ctx context.Context, addr string) error {
	s.initializeRoutes()
	s.l.Info("Server is running on", zap.String("port", addr))
	err := http.ListenAndServe(addr, s.Router)
	if err != nil {
		s.l.Fatal("can't start server:" + err.Error())
	}
	<-ctx.Done()

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	// TODO: stop server
	return nil
}

func (s *Server) initializeRoutes() {
	router := mux.NewRouter()
	router.Handle("/hello", HelloWorldHandler()).Methods("GET")
	router.Use(s.loggingMiddleware)
	s.Router = router
}

func HelloWorldHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//nolint:errcheck
		w.Write([]byte("Hello, World"))
	})
}
