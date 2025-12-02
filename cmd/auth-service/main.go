package main

import (
	"context"
	"os"
	"sync"
	"syscall"
	"time"

	"github.com/ayayaakasvin/auth-service/internal/config"
	httpserver "github.com/ayayaakasvin/auth-service/internal/http-server"
	"github.com/ayayaakasvin/auth-service/internal/logger"
	"github.com/ayayaakasvin/auth-service/internal/repository/postgresql"
	"github.com/ayayaakasvin/auth-service/internal/repository/valkey"
	"github.com/ayayaakasvin/auth-service/internal/services/jwtservice"
	goshutdownchannel "github.com/ayayaakasvin/go-shutdown-channel"

	_ "github.com/ayayaakasvin/auth-service/docs"
)

func main() {
	// core elements, main context used in lifecycle and wg to keep app alive
	mainCtx, cancel := context.WithCancel(context.Background())
	wg := new(sync.WaitGroup)
	wg.Add(1)

	// logger
	log := logger.SetupLogger("Auth-Service")

	// config
	cfg := config.MustLoadConfig()
	log.Infof("Configs retrieved: %v", cfg)

	s := goshutdownchannel.NewShutdown(mainCtx, cancel)
	s.Notify(os.Interrupt, syscall.SIGTERM)

	// dependencies
	cc := valkey.NewValkeyClient(cfg.ValkeyConfig, s)
	repo := postgresql.NewPostgreSQLConnection(cfg.PostgreSQLConfig, s)
	// cc := valkey.NewValkey_Mock()
	// repo := postgresql.NewPostgreSQL_Mock()
	jwtM := jwtservice.NewJWTManager(&cfg.JWTSecret)
	log.Info("Dependencies retrieved")

	app := httpserver.NewServerApp(s, &cfg.HTTPServer, log, repo, cc, jwtM)
	log.Info("New Application set up finished")

	go func() {
		defer wg.Done()

		<-s.Done()
		log.Println(s.Message())

		_, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		log.WithError(mainCtx.Err()).Info("Gracefully Shutdowning...")
	}()

	go app.Run()

	wg.Wait()
	log.WithError(mainCtx.Err()).Info("Graceful Shutdown completed")
}
