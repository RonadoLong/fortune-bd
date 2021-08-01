package delayer

import (
	"errors"
	"github.com/gomodule/redigo/redis"
	"time"
	"wq-fotune-backend/libs/logger"
)

// 键名
const (
	KeyJobPool       = "delayer:job_pool"
	PrefixJobBucket  = "delayer:job_bucket:"
	PrefixReadyQueue = "delayer:ready_queue:"
)

// 客户端结构
type Client struct {
	Conn     redis.Conn
	Host     string
	Database int
	Password string
}

// 初始化
func (p *Client) Init() error {
	// 创建连接
	if p.Conn == nil {
		db := 0
		conn, err := redis.Dial("tcp", p.Host, redis.DialPassword(p.Password), redis.DialDatabase(db))
		if nil != err {
			logger.Err(err)
			return err
		}
		p.Conn = conn
	}
	return nil
}

// 增加任务 单位是秒
func (p *Client) Push(message Message, delayTime int, readyMaxLifetime int) (bool, error) {
	// 参数验证
	if !message.Valid() {
		return false, errors.New("Invalid message")
	}
	// 执行事务
	_ = p.Conn.Send("MULTI")
	_ = p.Conn.Send("HMSET", PrefixJobBucket+message.ID, "topic", message.Topic, "body", message.Body)
	_ = p.Conn.Send("EXPIRE", PrefixJobBucket+message.ID, delayTime+readyMaxLifetime)
	_ = p.Conn.Send("ZADD", KeyJobPool, time.Now().Unix()+int64(delayTime), message.ID)
	values, err := redis.Values(p.Conn.Do("EXEC"))
	if err != nil {
		return false, err
	}
	// 事务结果处理
	v := values[0].(string)
	v1 := values[1].(int64)
	v2 := values[2].(int64)
	if v != "OK" || v1 == 0 || v2 == 0 {
		return false, nil
	}
	// 返回
	return true, nil
}

// 取出任务
func (p *Client) Pop(topic string) (*Message, error) {
	id, err := redis.String(p.Conn.Do("RPOP", PrefixReadyQueue+topic))
	if err != nil {
		return nil, err
	}
	result, err := redis.StringMap(p.Conn.Do("HGETALL", PrefixJobBucket+id))
	if err != nil {
		return nil, err
	}
	if result["topic"] == "" || result["body"] == "" {
		return nil, errors.New("Job bucket has expired or is incomplete")
	}
	_, _ = p.Conn.Do("DEL", PrefixJobBucket+id)
	msg := &Message{
		ID:    id,
		Topic: result["topic"],
		Body:  result["body"],
	}
	return msg, nil
}

// 阻塞取出任务
func (p *Client) BPop(topic string, timeout int) (*Message, error) {
	values, err := redis.Strings(p.Conn.Do("BRPOP", PrefixReadyQueue+topic, timeout))
	if err != nil {
		return nil, err
	}
	id := values[1]
	result, err := redis.StringMap(p.Conn.Do("HGETALL", PrefixJobBucket+id))
	if err != nil {
		return nil, err
	}
	if result["topic"] == "" || result["body"] == "" {
		return nil, errors.New("Job bucket has expired or is incomplete")
	}
	_, _ = p.Conn.Do("DEL", PrefixJobBucket+id)
	msg := &Message{
		ID:    id,
		Topic: result["topic"],
		Body:  result["body"],
	}
	return msg, nil
}

// 移除任务
func (p *Client) Remove(id string) (bool, error) {
	// 执行事务
	_ = p.Conn.Send("MULTI")
	_ = p.Conn.Send("ZREM", KeyJobPool, id)
	_ = p.Conn.Send("DEL", PrefixJobBucket+id)
	values, err := redis.Values(p.Conn.Do("EXEC"))
	if err != nil {
		return false, err
	}
	// 事务结果处理
	v := values[0].(int64)
	v1 := values[1].(int64)
	if v == 0 || v1 == 0 {
		return false, nil
	}
	// 返回
	return true, nil
}
