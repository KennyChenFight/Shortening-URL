package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/KennyChenFight/Shortening-URL/pkg/graceful"

	"github.com/KennyChenFight/golib/ratelimitlib"

	"github.com/KennyChenFight/Shortening-URL/pkg/repository"

	"github.com/KennyChenFight/Shortening-URL/pkg/lock"

	"github.com/KennyChenFight/Shortening-URL/pkg/middleware"

	"github.com/KennyChenFight/Shortening-URL/pkg/validation"
	"github.com/KennyChenFight/golib/loglib"
	"github.com/KennyChenFight/golib/migrationlib"
	"github.com/KennyChenFight/golib/pglib"
	"github.com/KennyChenFight/golib/redislib"
	"go.uber.org/zap"

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

type RedisConfig struct {
	URL string `long:"url" description:"redis url" env:"URL" required:"true"`
}

type SlideWindowRateLimiterConfig struct {
	Capacity int64         `long:"capacity" description:"max limit" env:"CAPACITY" required:"true"`
	Interval time.Duration `long:"interval" description:"window size" env:"INTERVAL" required:"true"`
}

type GinConfig struct {
	Port string `long:"port" description:"port" env:"PORT" default:":8080"`
	Mode string `long:"mode" description:"mode" env:"MODE" default:"debug"`
}

type Environment struct {
	GinConfig                    GinConfig                    `group:"gin" namespace:"Gin" env-namespace:"GIN"`
	PostgresConfig               PostgresConfig               `group:"postgres" namespace:"postgres" env-namespace:"POSTGRES"`
	RedisConfig                  RedisConfig                  `group:"redis" namespace:"redis" env-namespace:"REDIS"`
	SlideWindowRateLimiterConfig SlideWindowRateLimiterConfig `group:"slide-window-rate-limiter" namespace:"slide-window-rate-limiter" env-namespace:"SLIDE_WINDOW_RATE_LIMITER"`
	FQDN                         string                       `long:"fqdn" description:"fqdn" env:"FQDN" default:"localhost:8080"`
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

	redisClient, err := redislib.NewGORedisClient(redislib.GORedisConfig{URL: env.RedisConfig.URL}, nil)
	if err != nil {
		log.Fatalf("fail to init redis client:%v", err)
	}

	logger, err := loglib.NewProductionLogger()
	if err != nil {
		log.Fatalf("fail to init logger:%v", err)
	}

	urlDAO := dao.NewPGUrlDAO(logger, pgClient)
	cacheDAO := dao.NewRedisCacheDAO(logger, redisClient)

	bindingValidator, _ := binding.Validator.Engine().(*validator.Validate)
	CustomValidator, err := validation.NewValidationTranslator(bindingValidator, "en")
	if err != nil {
		log.Fatalf("fail to init validation translator:%v", err)
	}

	locker := lock.NewRedisLocker(logger, redisClient)

	rateLimiter := ratelimitlib.NewSlideWindowRateLimiter(redisClient, env.SlideWindowRateLimiterConfig.Capacity, env.SlideWindowRateLimiterConfig.Interval)

	mwe := middleware.NewMiddleware(logger, CustomValidator, rateLimiter)

	urlRepository := repository.NewURLRepository(logger, urlDAO, cacheDAO, locker)

	svc := service.NewService(&service.Config{FQDN: env.FQDN}, logger, urlRepository)

	gin.SetMode(env.GinConfig.Mode)

	graceful.Wrapper(logger, StartFunc(logger, server.NewHTTPServer(gin.Default(), env.GinConfig.Port, mwe, svc)))
}

func StartFunc(logger *loglib.Logger, server *http.Server) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		go func() {
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				logger.Error("http server listen error", zap.Error(err))
			}
		}()

		<-ctx.Done()

		ctx1, cancel1 := context.WithCancel(context.Background())

		go func() {
			logger.Info("shutdown http server...")
			server.Shutdown(ctx1)
			cancel1()
		}()

		<-ctx1.Done()
		logger.Info("http server existing")
		return nil
	}
}
