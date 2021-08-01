package env

var (
	developEnv = envConfig{
		runMode:   "dev",
		redisAddr: "192.168.3.30:16379",
		redisPWD:  "abc123++",
		etcdAddr:  "192.168.3.30:2379",
		dbDSN:     "root:abc123++@tcp(192.168.3.30:13306)/wq_fotune?charset=utf8mb4&parseTime=True&loc=Local",
		mongoAddr: "mongodb://wq:abc123@192.168.3.30:38888/ifortune",
	}
	releaseEnv = envConfig{
		runMode:   "dev",
		redisAddr: "192.168.3.30:16379",
		redisPWD:  "abc123++",
		etcdAddr:  "192.168.3.30:2379",
		dbDSN:     "root:abc123++@tcp(192.168.3.30:13306)/wq_fotune?charset=utf8mb4&parseTime=True&loc=Local",
		mongoAddr: "mongodb://wq:abc123@192.168.3.30:38888/ifortune",
	}
	proEnv = envConfig{
		runMode:   "dev",
		redisAddr: "192.168.3.30:16379",
		redisPWD:  "abc123++",
		etcdAddr:  "192.168.3.30:2379",
		dbDSN:     "root:abc123++@tcp(192.168.3.30:13306)/wq_fotune?charset=utf8mb4&parseTime=True&loc=Local",
		mongoAddr: "mongodb://wq:abc123@192.168.3.30:38888/ifortune",
	}
)
