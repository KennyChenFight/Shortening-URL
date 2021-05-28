package main

import (
	"context"
	"log"
	"os"

	"github.com/KennyChenFight/Shortening-URL/pkg/dao"
	"github.com/KennyChenFight/Shortening-URL/pkg/graceful"
	"github.com/KennyChenFight/Shortening-URL/pkg/job"
	"github.com/KennyChenFight/golib/loglib"
	"github.com/KennyChenFight/golib/pglib"
	"github.com/KennyChenFight/randstr"
	"github.com/jessevdk/go-flags"
)

type PostgresConfig struct {
	URL              string `long:"url" description:"database url" env:"URL" required:"true"`
	DebugMode        bool   `long:"debug-mode" description:"database debug mode" env:"DEBUG_MODE"`
	PoolSize         int    `long:"pool-size" description:"database pool size" env:"POOL_SIZE" default:"100"`
	MigrationFileDir string `long:"migration-file-dir" description:"migration file dir" env:"MIGRATION_FILE_DIR" default:"file://migrations"`
}

type GenerateKeyConfig struct {
	EveryKeyNumber int `long:"every-key-number" description:"every key number" env:"EVERY_KEY_NUMBER" default:"1000000"`
}

type ExpireURLConfig struct {
	Number int `long:"number" description:"expire url number" env:"NUMBER" default:"1000"`
}

type Environment struct {
	PostgresConfig    PostgresConfig    `group:"postgres" namespace:"postgres" env-namespace:"POSTGRES"`
	GenerateKeyConfig GenerateKeyConfig `group:"generate-key" namespace:"generate-key" env-namespace:"GENERATE_KEY"`
	ExpireURLConfig   ExpireURLConfig   `group:"expire-url" namespace:"expire-url" env-namespace:"EXPIRE_URL"`
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

	pgClient, err := pglib.NewDefaultGOPGClient(pglib.GOPGConfig{
		URL:       env.PostgresConfig.URL,
		DebugMode: env.PostgresConfig.DebugMode,
		PoolSize:  env.PostgresConfig.PoolSize,
	})
	if err != nil {
		log.Fatalf("fail to init postgres client:%v", err)
	}

	logger, err := loglib.NewProductionLogger()
	if err != nil {
		log.Fatalf("fail to init logger:%v", err)
	}

	randomStrGenerator := randstr.NewFastGenerator(randstr.CharSetEnglishAlphabet)

	keyDAO := dao.NewPGKeyDAO(logger, pgClient, randomStrGenerator)
	urlDAO := dao.NewPGUrlDAO(logger, pgClient)

	generateKeyJob := job.NewGenerateKeyJob(job.GenerateKeyJobConfig{Name: "GenerateKeyJob", TimerFormat: "5 4 * * sun", EveryKeyNumber: env.GenerateKeyConfig.EveryKeyNumber}, keyDAO)
	expireURLJob := job.NewExpiredURLJob(job.ExpiredURLJobConfig{Name: "ExpiredURLJob", TimerFormat: "5 4 * * sun", ExpireURLNumber: env.ExpireURLConfig.Number}, urlDAO)
	jobs := []job.Job{generateKeyJob, expireURLJob}
	manager := job.NewManager(jobs, logger)

	graceful.Wrapper(logger, StartFunc(logger, manager))
}

func StartFunc(logger *loglib.Logger, manager *job.Manager) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		go func() {
			manager.Start()
			logger.Info("start cronjob...")
		}()

		<-ctx.Done()

		ctx1, cancel1 := context.WithCancel(context.Background())

		go func() {
			logger.Info("shutdown cronjob manager...")
			manager.Stop()
			cancel1()
		}()

		<-ctx1.Done()
		logger.Info("cronjob manager existing")
		return nil
	}
}
