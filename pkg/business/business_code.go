package business

const (
	// common
	Unknown           = 1000
	NotFound          = 1001
	Validation        = 1002
	Internal          = 1003
	ServerUnavailable = 1004
	MethodNowAllowed  = 1005
	PathNotFound      = 1006
	TooManyRequest    = 1007

	// postgres
	PostgresInternalError = 1100

	// redis
	RedisInternalError = 1200

	// lock
	AcquireLockURLResourceError = 1300
)
