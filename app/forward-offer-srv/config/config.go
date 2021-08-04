package config

import (
	"github.com/spf13/viper"
	"os"
	"wq-fotune-backend/libs/logger"
)

const (
	mysqlHostName = "mysql.host"
	redisHostName = "redis.host"
	redisPassword = "redis.password"
	etcdAddr      = "etcd.addr"
)

const (
	SRV_NAME  = "forward-offer-srv"
	JaegerUrl = ""
)

var Config *conf

type conf struct {
	MysqlHost string
	RunMode   string
	RedisHost string
	RedisPass string
	EtcdAddr  string
}

func Init(configPath string) {
	var err error
	v := viper.New()
	v.SetConfigFile(configPath)

	err = v.ReadInConfig()
	if err != nil {
		logger.Errorf("%s", err.Error())
		os.Exit(-1)
	}
	logger.Infof("当前配置文件路径是 ======= %s", configPath)
	conf := conf{}
	conf.RunMode = v.GetString("runMode")
	logger.Infof("当前环境是 ======= %s", conf.RunMode)

	conf.MysqlHost = v.GetString(mysqlHostName)
	conf.RedisHost = v.GetString(redisHostName)
	conf.RedisPass = v.GetString(redisPassword)
	conf.EtcdAddr = v.GetString(etcdAddr)
	Config = &conf
}
