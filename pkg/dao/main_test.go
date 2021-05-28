package dao

import (
	"os"
	"strconv"
	"testing"

	"github.com/KennyChenFight/golib/migrationlib"
	"github.com/KennyChenFight/golib/pglib"
	"github.com/KennyChenFight/golib/redislib"
	"github.com/golang-migrate/migrate/v4"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDao(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Dao Suite")
}

var testRedisClient *redislib.GORedisClient
var testPGClient *pglib.GOPGClient

var _ = BeforeSuite(func() {
	testPGClient = setupTestPG()
	testRedisClient = setupTestRedis()
})

var _ = AfterSuite(func() {
	testPGClient.Close()
	testRedisClient.Close()
})

func setupTestPG() *pglib.GOPGClient {
	pgURL := os.Getenv("POSTGRES_URL")
	if pgURL == "" {
		panic("should setup postgres url")
	}

	pgDebugMode, _ := strconv.ParseBool(os.Getenv("pgDebugMode"))
	pgClient, err := pglib.NewDefaultGOPGClient(pglib.GOPGConfig{
		URL:       pgURL,
		DebugMode: pgDebugMode,
		PoolSize:  10,
	})
	Expect(err).To(BeNil())
	Expect(pgClient).NotTo(BeNil())

	migrationFileDir := os.Getenv("POSTGRES_MIGRATION_FILE_DIR")
	if migrationFileDir == "" {
		migrationFileDir = "file://../../migrations"
	}

	migration := migrationlib.NewMigrateLib(migrationlib.Config{
		DatabaseDriver: migrationlib.PostgresDriver,
		DatabaseURL:    pgURL,
		SourceDriver:   migrationlib.FileDriver,
		SourceURL:      migrationFileDir,
		TableName:      "migrate_version",
	})
	err = migration.Up()
	if err == migrate.ErrNoChange {
		err = nil
	}
	Expect(err).To(BeNil())
	return pgClient
}

func setupTestRedis() *redislib.GORedisClient {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		panic("should setup redis url")
	}

	redisClient, err := redislib.NewGORedisClient(redislib.GORedisConfig{URL: redisURL}, nil)
	Expect(err).To(BeNil())
	Expect(redisClient).NotTo(BeNil())

	return redisClient
}
