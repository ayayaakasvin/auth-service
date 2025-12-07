package httpserver

import (
	"context"
	"net/http"
	"time"

	"github.com/ayayaakasvin/auth-service/internal/config"
	"github.com/ayayaakasvin/auth-service/internal/http-server/handlers"
	"github.com/ayayaakasvin/auth-service/internal/http-server/middlewares"
	"github.com/ayayaakasvin/auth-service/internal/models/core"
	"github.com/ayayaakasvin/auth-service/internal/services/jwtservice"
	goshutdownchannel "github.com/ayayaakasvin/go-shutdown-channel"
	"github.com/ayayaakasvin/lightmux"
	"github.com/sirupsen/logrus"

	httpSwagger "github.com/swaggo/http-swagger"
)

type ServerApp struct {
	server *http.Server

	lmux *lightmux.LightMux
	s    *goshutdownchannel.Shutdown

	repo  core.Repository
	cache core.Cache
	jwtM  *jwtservice.JWTService

	cfg *config.HTTPServer

	logger *logrus.Logger
}

func NewServerApp(
	s *goshutdownchannel.Shutdown,
	httpcfg *config.HTTPServer,
	corscfg *config.CorsConfig,
	logger *logrus.Logger,
	repo core.Repository,
	cache core.Cache,
	jwtM *jwtservice.JWTService,
) *ServerApp {
	return &ServerApp{
		cfg:    httpcfg,
		logger: logger,
		repo:   repo,
		cache:  cache,
		s:      s,
		jwtM:   jwtM,
	}
}

func (s *ServerApp) Run() {
	s.setupServer()

	s.setupLightMux()

	s.startServer()
}

func (s *ServerApp) startServer() {
	s.logger.Infof("Server has been started on port: %s", s.cfg.Address)
	s.logger.Infof("Available handlers:\n")

	s.lmux.PrintMiddlewareInfo()
	s.lmux.PrintRoutes()

	go printServerStatus(s.s.Context(), s.logger)

	// RunTLS can be run when server is hosted on domain, acts as seperate service of file storing, for my project, id chose to encapsulate servers under one docker-compose and make nginx-gateaway for my api like auth, file, user service
	// if err := s.lmux.RunTLS(s.cfg.TLS.CertFile, s.cfg.TLS.KeyFile); err != nil {
	if err := s.lmux.RunContext(s.s.Context()); err != nil {
		s.logger.Fatalf("Server exited with error: %v", err)
	}
}

// setuping server by pointer, so we dont have to return any value
func (s *ServerApp) setupServer() {
	if s.server == nil {
		// s.logger.Warn("Server is nil, creating a new server pointer")
		s.server = &http.Server{}
	}

	s.server.Addr = s.cfg.Address
	s.server.IdleTimeout = s.cfg.IdleTimeout
	s.server.ReadTimeout = s.cfg.Timeout
	s.server.WriteTimeout = s.cfg.Timeout

	s.logger.Info("Server has been set up")
}

func (s *ServerApp) setupLightMux() {
	s.lmux = lightmux.NewLightMux(s.server)

	mws := middlewares.NewHTTPMiddlewares(s.logger, config.CorsConfig{}, s.cache,s.jwtM)
	hndlrs := handlers.NewHTTPHandlers(s.repo, s.cache, s.logger, s.jwtM)

	s.lmux.Use(mws.RecoverMiddleware, mws.LoggerMiddleware, mws.CORSMiddleware)

	s.lmux.NewRoute("/ping").Handle(http.MethodGet, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("pong"))
	})
	// s.lmux.NewRoute("/panic").Handle(http.MethodGet, PanicHandler())

	apiGroup := s.lmux.NewGroup("/api")

	apiGroup.NewRoute("/login", mws.RateLimitLoginMiddleware).Handle(http.MethodPost, hndlrs.LogIn())
	apiGroup.NewRoute("/register", mws.RateLimitRegisterMiddleware).Handle(http.MethodPost, hndlrs.Register())
	apiGroup.NewRoute("/logout", mws.JWTAuthMiddleware).Handle(http.MethodDelete, hndlrs.LogOut())
	apiGroup.NewRoute("/refresh").Handle(http.MethodPost, hndlrs.RefreshTheToken())

	apiGroup.NewRoute("/public/user").Handle(http.MethodGet, hndlrs.PublicUserInfo())
	apiGroup.NewRoute("/me", mws.JWTAuthMiddleware).Handle(http.MethodGet, hndlrs.PrivateUserInfo())

	s.lmux.Mux().HandleFunc("/swagger/", httpSwagger.WrapHandler)

	s.logger.Info("LightMux has been set up")
}

func printServerStatus(ctx context.Context, log *logrus.Logger) {
	ticker := time.NewTicker(time.Minute * 1)

	for {
		select {
		case <-ticker.C:
			log.Info("Server is alive...")
		case <-ctx.Done():
			return
		}
	}
}

// Used for recover test
func PanicHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		panic("ambatubas")
	}
}