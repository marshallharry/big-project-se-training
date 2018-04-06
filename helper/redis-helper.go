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
	redisPool = newRedis("127.0.0.1:6379")
	subscribe("mh0318-nsq-visitor-incr", "ch1", handlerIncrementVisitor())
	subscribe("mh0318-nsq-visitor-set", "ch1", handlerSetVisitor())
}

func newRedis(address string) *redigo.Pool {
	return &redigo.Pool{
		MaxIdle:     1,
		IdleTimeout: 10 * time.Second,
		Dial:        func() (redigo.Conn, error) { return redigo.Dial("tcp", address) },
	}
}

func setRedis(key string, value interface{}) error {
	con := redisPool.Get()
	defer con.Close()

	_, err := con.Do("SET", key, value)
	return err
}

// GetRedis handler
func GetRedis(key string) (string, error) {
	con := redisPool.Get()
	defer con.Close()

	return redigo.String(con.Do("GET", key))
}

func increment(key string) error {
	con := redisPool.Get()
	defer con.Close()

	_, err := con.Do("INCR", key)
	return err
}

func handlerIncrementVisitor() nsq.HandlerFunc {
	return nsq.HandlerFunc(func(message *nsq.Message) error {
		increment("mh0318-redis-visitor")
		return nil
	})
}

func handlerSetVisitor() nsq.HandlerFunc {
	return nsq.HandlerFunc(func(message *nsq.Message) error {
		setRedis("mh0318-redis-visitor", 1)
		return nil
	})
}
