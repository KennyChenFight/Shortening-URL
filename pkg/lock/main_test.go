package lock

import (
	"os"
	"testing"

	"github.com/KennyChenFight/golib/redislib"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestLock(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Lock Suite")
}

var testRedisClient *redislib.GORedisClient

var _ = BeforeSuite(func() {
	testRedisClient = setupTestRedis()
})

var _ = AfterSuite(func() {
	testRedisClient.Close()
})

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
