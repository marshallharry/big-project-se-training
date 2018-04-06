package helper

import (
	"time"

	"github.com/bitly/go-nsq"
	redigo "github.com/garyburd/redigo/redis"
)

var (
	redisPool *redigo.Pool
	subscribe = Subscribe
)

func init() {
	redisPool = NewRedis("127.0.0.1:6379")
	subscribe("mh0318-nsq-visitor-incr", "ch1", HandlerIncrementVisitor())
	subscribe("mh0318-nsq-visitor-set", "ch1", HandlerSetVisitor())
}

func NewRedis(address string) *redigo.Pool {
	return &redigo.Pool{
		MaxIdle:     1,
		IdleTimeout: 10 * time.Second,
		Dial:        func() (redigo.Conn, error) { return redigo.Dial("tcp", address) },
	}
}

func SetRedis(key string, value interface{}) error {
	con := redisPool.Get()
	defer con.Close()

	_, err := con.Do("SET", key, value)
	return err
}

func GetRedis(key string) (string, error) {
	con := redisPool.Get()
	defer con.Close()

	return redigo.String(con.Do("GET", key))
}

func Increment(key string) error {
	con := redisPool.Get()
	defer con.Close()

	_, err := con.Do("INCR", key)
	return err
}

func HandlerIncrementVisitor() nsq.HandlerFunc {
	return nsq.HandlerFunc(func(message *nsq.Message) error {
		Increment("mh0318-redis-visitor")
		return nil
	})
}

func HandlerSetVisitor() nsq.HandlerFunc {
	return nsq.HandlerFunc(func(message *nsq.Message) error {
		SetRedis("mh0318-redis-visitor", 1)
		return nil
	})
}
