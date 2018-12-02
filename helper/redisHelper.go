package helper

import (
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"os"
	"time"
)

func CreateRedisPool() *redis.Pool {
	redisAddr := os.Getenv("REDIS_ADDR")
	redisPwd := os.Getenv("REDIS_PWD")
	return GetRedisPool(redisAddr, redisPwd)
}

func GetRedisPool(addr, pwd string) *redis.Pool {
	return GetRedisDBPool(addr, pwd, 0)
}

//EX seconds − 设置指定的到期时间(以秒为单位)。PX milliseconds - 设置指定的到期时间(以毫秒为单位)。
//NX - 仅在键不存在时设置键。
//XX - 只有在键已存在时才设置。

func GetRedisDBPool(addr, pwd string, db int) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     6,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			opts := []redis.DialOption{
				redis.DialConnectTimeout(5 * time.Second),
				redis.DialReadTimeout(2 * time.Second),
				redis.DialWriteTimeout(2 * time.Second),
				redis.DialDatabase(db),
				redis.DialPassword(pwd),
			}
			c, err := redis.Dial("tcp", addr, opts...)
			if err != nil {
				fmt.Printf("connect %s error %s", addr, err.Error())
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func GetRedisMaxActPool(addr, pwd string) *redis.Pool {
	return GetRedisDBMaxActPool(addr, pwd, 0)
}

func GetRedisDBMaxActPool(addr, pwd string, db int) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     6,
		IdleTimeout: 240 * time.Second,
		MaxActive:   200,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			opts := []redis.DialOption{
				redis.DialConnectTimeout(5 * time.Second),
				redis.DialReadTimeout(2 * time.Second),
				redis.DialWriteTimeout(2 * time.Second),
				redis.DialDatabase(db),
				redis.DialPassword(pwd),
			}
			c, err := redis.Dial("tcp", addr, opts...)
			if err != nil {
				fmt.Printf("connect %s error %s", addr, err.Error())
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

type Mutex struct {
	key    string
	value  string
	expire int
	pool   *redis.Pool
}

// redis分布式锁客户端，在db1保存
//func NewMutex(key, value string, expire int) *Mutex {
//	mutex := &Mutex{}
//	conf := smconfig.GetDBConfig()
//	mutex.key = key     // key为了分辨是否已经锁
//	mutex.value = value // value 为了分辨谁锁的
//	mutex.expire = expire // 毫秒
//	mutex.pool = GetRedisDBPool(conf.RedisRoomUserAddr, conf.RedisRoomUserPwd, 1)
//	return mutex
//}

const (
	LOCK_SUCCESS        = "OK"
	DEFAULT_EXPIRE_TIME = 180000 // 3分钟
)

var (
	ErrLockHeld = errors.New("ErrLockHeld")
)

func (m *Mutex) Lock() (bool, error) {
	redisDB := m.pool.Get()
	defer redisDB.Close()

	expireTime := DEFAULT_EXPIRE_TIME
	if m.expire != 0 {
		expireTime = m.expire
	}
	res, err := redis.String(redisDB.Do("SET", m.key, m.value, "NX", "PX", expireTime))
	if err == redis.ErrNil {
		// 返回ErrLockHeld表示锁已经被其他服务拿到了
		return false, ErrLockHeld
	}

	if err != nil {
		//logger.Errorf(nil, "redis Lock err %v", err)
		return false, err
	}
	if res == LOCK_SUCCESS {
		return true, nil
	}

	return false, nil
}

func (m *Mutex) UnLock() bool {
	redisDB := m.pool.Get()
	defer redisDB.Close()

	res, err := redis.Int(delScript.Do(redisDB, m.key, m.value))
	if err != nil {
		fmt.Printf("redis UnLock err %v", err)
		return false
	}

	if res == 1 {
		return true
	}

	return false
}

func (m *Mutex) Touch() bool {
	redisDB := m.pool.Get()
	defer redisDB.Close()

	res, err := redis.String(touchScript.Do(redisDB, m.key, m.value, DEFAULT_EXPIRE_TIME))
	if err != nil {
		fmt.Printf("redis Touch err %v", err)
		return false
	}

	if res == LOCK_SUCCESS {
		return true
	}

	return false
}

func (m *Mutex) TryLockTimes(waitIndex int, waitTime int) (bool, error) {
	for i := 0; i < waitIndex; i++ {
		res, err := m.Lock()
		if res {
			return res, nil
		}
		if err == ErrLockHeld {
			//logger.Debugf(nil, "Lock waiting... %v time", i)
			time.Sleep(time.Duration(waitTime) * time.Millisecond)
			continue
		}
		if err != nil {
			return false, err
		}
	}
	return false, nil
}

var delScript = redis.NewScript(1, `
if redis.call("get", KEYS[1]) == ARGV[1] then
	return redis.call("del", KEYS[1])
else
	return 0
end`)

var touchScript = redis.NewScript(1, `
if redis.call("get", KEYS[1]) == ARGV[1] then
	return redis.call("set", KEYS[1], ARGV[1], "xx", "px", ARGV[2])
else
	return "ERR"
end`)
