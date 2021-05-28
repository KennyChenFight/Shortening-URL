package main

import (
	"context"
	"github.com/KennyChenFight/Shortening-URL/pkg/validation"
	"github.com/KennyChenFight/golib/loglib"
	"github.com/KennyChenFight/golib/migrationlib"
	"github.com/KennyChenFight/golib/pglib"
	"github.com/KennyChenFight/randstr"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang-migrate/migrate/v4"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"

	"github.com/KennyChenFight/Shortening-URL/pkg/dao"
	"github.com/KennyChenFight/Shortening-URL/pkg/server"
	"github.com/KennyChenFight/Shortening-URL/pkg/service"
	"github.com/gin-gonic/gin"
	"github.com/jessevdk/go-flags"
)

type PostgresConfig struct {
	URL              string `long:"url" description:"database url" env:"URL" required:"true"`
	MigrationFileDir string `long:"migration-file-dir" description:"migration file dir" env:"MIGRATION_FILE_DIR" default:"file://migrations"`
}

type GinConfig struct {
	Port string `long:"port" description:"port" env:"PORT" default:":8080"`
	Mode string `long:"mode" description:"mode" env:"MODE" default:"debug"`
}

type Environment struct {
	GinConfig      GinConfig      `group:"gin" namespace:"Gin" env-namespace:"GIN"`
	PostgresConfig PostgresConfig `group:"postgres" namespace:"postgres" env-namespace:"POSTGRES"`
	FQDN           string         `long:"fqdn" description:"fqdn" env:"FQDN" default:"localhost:8080"`
}

func main() {
	var env Environment
	parser := flags.NewParser(&env, flags.Default)
	if _, err := parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

	migration := migrationlib.NewMigrateLib(migrationlib.Config{
		DatabaseDriver: migrationlib.PostgresDriver,
		DatabaseURL:    env.PostgresConfig.URL,
		SourceDriver:   migrationlib.FileDriver,
		SourceURL:      env.PostgresConfig.MigrationFileDir,
		TableName:      "migrate_version",
	})
	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("fail to run database migration fail:%v", err)
	}

	pgClient, err := pglib.NewDefaultGOPGClient(pglib.GOPGConfig{
		URL:       env.PostgresConfig.URL,
		DebugMode: false,
		PoolSize:  10,
	})
	if err != nil {
		log.Fatalf("fail to init postgres client:%v", err)
	}

	logger, err := loglib.NewProductionLogger()
	if err != nil {
		log.Fatalf("fail to init logger:%v", err)
	}

	randomStrGenerator := randstr.NewFastGenerator(randstr.CharSetBase62)

	urlDAO := dao.NewPGUrlDAO(logger, pgClient, randomStrGenerator)

	bindingValidator, _ := binding.Validator.Engine().(*validator.Validate)
	CustomValidator, err := validation.NewValidationTranslator(bindingValidator, "en")
	if err != nil {
		log.Fatalf("fail to init validation translator:%v", err)
	}

	svc := service.NewService(&service.Config{FQDN: env.FQDN}, logger, urlDAO, CustomValidator)

	gin.SetMode(env.GinConfig.Mode)
	StartRun(logger, server.NewHTTPServer(gin.Default(), env.GinConfig.Port, svc))
}

func StartRun(logger *loglib.Logger, server *http.Server) {
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("listen error", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down http server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exiting")
}
