package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/dormitory-life/gateway/internal/config"
)

type ServerConfig struct {
	Config         config.ServerConfig
	AuthServiceUrl string
	CoreServiceUrl string
	Logger         *slog.Logger
	JWTSecret      string
}

type Server struct {
	server    http.Server
	logger    *slog.Logger
	authProxy *httputil.ReverseProxy
	coreProxy *httputil.ReverseProxy
	jwtSecret []byte
}

func New(cfg ServerConfig) *Server {
	s := new(Server)
	s.server.Addr = fmt.Sprintf(":%d", cfg.Config.Port)

	s.authProxy = s.createProxy(cfg.AuthServiceUrl, "/auth")
	s.coreProxy = s.createProxy(cfg.CoreServiceUrl, "/core")

	s.logger = cfg.Logger
	s.jwtSecret = []byte(cfg.JWTSecret)

	mux := s.setupRouter()
	s.server.Handler = s.loggingMiddleware(mux)

	return s
}

func (s *Server) createProxy(targetUrl string, prefix string) *httputil.ReverseProxy {
	target, _ := url.Parse(targetUrl)

	proxy := httputil.NewSingleHostReverseProxy(target)

	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.Host = target.Host
	}

	return proxy
}

func (s *Server) setupRouter() *http.ServeMux {
	mux := http.NewServeMux()

	s.setupPublicRoutes(mux)
	s.setupProtectedRoutes(mux)

	return mux
}

func (s *Server) setupPublicRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /auth/ping", s.pingHandler)
	mux.HandleFunc("POST /auth/register", s.authProxy.ServeHTTP)
	mux.HandleFunc("POST /auth/login", s.authProxy.ServeHTTP)
	mux.HandleFunc("POST /auth/refresh", s.authProxy.ServeHTTP)
}

func (s *Server) setupProtectedRoutes(mux *http.ServeMux) {
	mux.Handle("/auth/", s.authMiddleware(
		http.StripPrefix("/auth", s.authProxy),
	))

	mux.Handle("/core/", s.authMiddleware(
		http.StripPrefix("/core", s.coreProxy),
	))
}

func (s *Server) Start() error {
	s.logger.Debug("gateway server started", slog.String("address", s.server.Addr))
	return s.server.ListenAndServe()
}
