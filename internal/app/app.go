package app

import (
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"todolist_api/config"
	v1 "todolist_api/internal/api/v1"
	"todolist_api/internal/repo"
	"todolist_api/internal/service"
	"todolist_api/pkg/hasher"
	"todolist_api/pkg/httpserver"
	"todolist_api/pkg/postgres"
	"todolist_api/pkg/validator"
)

//	@title			Api for tasks
//	@version		1.0
//	@description	Api for tasks. Include create, update, delete tasks

//	@host		localhost:8080
//	@BasePath	/

//	@securityDefinitions.apikey	JWT
//	@in							header
//	@name						Authorization
//	@description				JWT token

func Run() {
	// config
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}
	// set up json logger
	setLogger(cfg.Log.Level, cfg.Log.Output)

	// postgresql database
	pg, err := postgres.NewPG(cfg.PG.Url, postgres.MaxPoolSize(cfg.PG.MaxPoolSize))
	if err != nil {
		log.Fatalf("Initializing postgres error: %s", err)
	}
	defer pg.Close()

	d := &service.ServicesDependencies{
		Repos:    repo.NewRepositories(pg),
		Hasher:   hasher.NewHasher(cfg.Hasher.Secret),
		SignKey:  cfg.JWT.SignKey,
		TokenTTL: cfg.JWT.TokenTTL,
	}
	services := service.NewServices(d)

	// validator for incoming messages
	v, err := validator.NewValidator()
	if err != nil {
		log.Fatalf("Initializing handler validator error: %s", err)
	}

	// handler
	handler := echo.New()
	handler.Validator = v
	v1.LoggingMiddleware(handler, cfg.Log.Output)
	v1.NewRouter(handler, services)

	// http server
	httpServer := httpserver.NewServer(handler, httpserver.Port(cfg.HTTP.Port))

	log.Infof("App started! Listening port %s", cfg.HTTP.Port)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Info("app run, signal " + s.String())

	case err = <-httpServer.Notify():
		log.Errorf("/app/run http server notify error: %s", err)
	}
	// graceful shutdown
	if err = httpServer.Shutdown(); err != nil {
		log.Errorf("/app/run http server shutdown error: %s", err)
	}

	log.Infof("App shutdown with exit code 0")
}

// loading environment params from .env
func init() {
	if _, ok := os.LookupEnv("HTTP_PORT"); !ok {
		if err := godotenv.Load(); err != nil {
			log.Fatalf("load env file error: %s", err)
		}
	}
}
