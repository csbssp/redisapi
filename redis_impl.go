package redisapi

import (
	"reflect"
	"strings"
	"time"

	"github.com/garyburd/redigo/redis"
)

type RedisClient struct {
	pool *redis.Pool
	addr string
}

func (rc RedisClient) Exists(key string) bool {
	conn := rc.connectInit()
	defer conn.Close()
	v, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false
	}
	return v
}

func (rc RedisClient) Lpush(key string, value interface{}) error {
	conn := rc.connectInit()
	defer conn.Close()

	if reflect.TypeOf(value).Kind() == reflect.Slice {
		s := reflect.ValueOf(value)
		values := make([]interface{}, s.Len()+1)
		values[0] = key
		for i := 1; i <= s.Len(); i++ {
			values[i] = s.Index(i - 1).Interface()
		}
		_, err := conn.Do("LPUSH", values...)
		return err
	} else {
		_, err := conn.Do("LPUSH", key, value)
		return err
	}
}

func (rc RedisClient) Rpush(key string, value interface{}) error {
	conn := rc.connectInit()
	defer conn.Close()

	if reflect.TypeOf(value).Kind() == reflect.Slice {
		s := reflect.ValueOf(value)
		values := make([]interface{}, s.Len()+1)
		values[0] = key
		for i := 1; i <= s.Len(); i++ {
			values[i] = s.Index(i - 1).Interface()
		}
		_, err := conn.Do("RPUSH", values...)
		return err
	} else {
		_, err := conn.Do("RPUSH", key, value)
		return err
	}
}

func (rc RedisClient) Lrange(key string, start, end int) ([]interface{}, error) {
	conn := rc.connectInit()
	defer conn.Close()

	v, err := conn.Do("LRANGE", key, start, end)
	return v.([]interface{}), err
}
func (rc RedisClient) Rpop(key string) (interface{}, error) {
	conn := rc.connectInit()
	defer conn.Close()

	value, err := conn.Do("RPOP", key)
	if err != nil {
		return nil, err
	}
	if value == nil {
		return nil, nil
	}
	return value, err
}

func (this RedisClient) Lset(key string, index int, value interface{}) error {
	conn := this.connectInit()
	defer conn.Close()

	_, err := conn.Do("LSET", key, index, value)
	return err
}

func (this RedisClient) Ltrim(key string, start, end int) error {
	conn := this.connectInit()
	defer conn.Close()

	_, err := conn.Do("LTRIM", key, start, end)
	return err
}

func (rc RedisClient) Brpop(key string, timeoutSecs int) (interface{}, error) {
	conn := rc.connectInit()
	defer conn.Close()

	var val interface{}
	var err error
	if timeoutSecs < 0 {
		val, err = conn.Do("BRPOP", key, 0)
	} else {
		val, err = conn.Do("BRPOP", key, timeoutSecs)
	}
	values, err := redis.Values(val, err)
	if err != nil {
		return nil, err
	}
	return string(values[1].([]byte)), err
}

func (rc RedisClient) Lrem(key string, value interface{}, remType int) error {
	conn := rc.connectInit()
	defer conn.Close()

	_, err := conn.Do("LREM", key, remType, value)
	return err
}

func (rc RedisClient) Set(key string, value []byte) error {
	conn := rc.connectInit()
	defer conn.Close()

	_, err := conn.Do("SET", key, value)
	return err
}

func (rc RedisClient) Get(key string) ([]byte, error) {
	conn := rc.connectInit()
	defer conn.Close()

	v, err := conn.Do("GET", key)
	if err != nil || v == nil {
		return nil, err
	}
	return v.([]byte), err
}

func (rc RedisClient) Delete(key string) error {
	conn := rc.connectInit()
	defer conn.Close()

	_, err := conn.Do("DEL", key)
	return err
}

func (rc RedisClient) Incr(key string, step uint64) (int64, error) {
	conn := rc.connectInit()
	defer conn.Close()

	value, err := conn.Do("INCRBY", key, step)
	if err != nil {
		return 0, nil
	}
	return value.(int64), err
}

func (rc RedisClient) Decr(key string, step uint64) (int64, error) {
	conn := rc.connectInit()
	defer conn.Close()

	value, err := conn.Do("DECRBY", key, step)
	if err != nil {
		return 0, nil
	}
	return value.(int64), err
}

func (rc RedisClient) MultiGet(keys []interface{}) ([]interface{}, error) {
	conn := rc.connectInit()
	defer conn.Close()

	v, err := conn.Do("MGET", keys...)
	return v.([]interface{}), err
}

func (rc RedisClient) MultiSet(kvMap map[string][]byte) error {
	conn := rc.connectInit()
	defer conn.Close()

	var values []interface{}
	for key, value := range kvMap {
		values = append(values, key)
		values = append(values, value)
	}

	_, err := conn.Do("MSET", values...)
	return err
}

func (rc RedisClient) ClearAll() error {
	conn := rc.connectInit()
	defer conn.Close()

	_, err := conn.Do("FLUSHALL")
	return err
}

// order set begin
func (this RedisClient) Zadd(key string, score int, value interface{}) error {
	conn := this.connectInit()
	defer conn.Close()

	_, err := conn.Do("ZADD", key, score, value)
	return err
}

func (this RedisClient) Zrem(key string, value interface{}) error {
	conn := this.connectInit()
	defer conn.Close()

	_, err := conn.Do("ZREM", key, value)
	return err
}

func (this RedisClient) ZRemRangeByRank(key string, start, end int) error {
	conn := this.connectInit()
	defer conn.Close()

	_, err := conn.Do("ZREMRANGEBYRANK", key, start, end)
	return err
}

func (this RedisClient) Zcard(key string) (int, error) {
	conn := this.connectInit()
	defer conn.Close()

	v, err := redis.Int(conn.Do("ZCARD", key))
	return int(v), err
}

func (this RedisClient) ZRrank(key string, value interface{}) (int, error) {
	conn := this.connectInit()
	defer conn.Close()

	v, err := redis.Int(conn.Do("ZRANK", key, value))
	return int(v), err
}

func (this RedisClient) ZRrange(key string, begin int, end int) (scoreStructList []ScoreStruct, err error) {
	conn := this.connectInit()
	defer conn.Close()

	v, err := redis.Values(conn.Do("ZRANGE", key, begin, end, "WITHSCORES"))
	if err != nil {
		return nil, err
	}
	length := len(v)
	scoreStructList = make([]ScoreStruct, length/2)
	for i := 0; i < length/2; i++ {
		scoreStructList[i].Member = v[i*2]
		scoreStructList[i].Score = v[i*2+1]
	}
	return
}

func (this RedisClient) ZRevRrank(key string, value interface{}) (int, error) {
	conn := this.connectInit()
	defer conn.Close()

	v, err := redis.Int(conn.Do("ZREVRANK", key, value))
	return int(v), err
}

func (this RedisClient) ZRevRrange(key string, begin int, end int) (scoreStructList []ScoreStruct, err error) {
	conn := this.connectInit()
	defer conn.Close()

	v, err := redis.Values(conn.Do("ZREVRANGE", key, begin, end, "WITHSCORES"))
	if err != nil {
		return nil, err
	}
	length := len(v)
	scoreStructList = make([]ScoreStruct, length/2)
	for i := 0; i < length/2; i++ {
		scoreStructList[i].Member = v[i*2]
		scoreStructList[i].Score = v[i*2+1]
	}
	return
}

// order set end

func (rc RedisClient) Pub(key string, value interface{}) error {
	conn := rc.connectInit()
	defer conn.Close()

	_, err := conn.Do("PUBLISH", key, value)
	return err
}

func (rc RedisClient) Sub(keys ...string) ([]string, error) {
	conn := rc.connectInit()
	defer conn.Close()

	v, err := conn.Do("SUBSCRIBE", keys)
	value := string(v.([]interface{})[1].([]byte))
	values := strings.Split(value[1:len(value)-1], " ")
	return values, err
}

func (rc RedisClient) UnSub(keys ...string) error {
	conn := rc.connectInit()
	defer conn.Close()

	_, err := conn.Do("UNSUBSCRIBE", keys)
	return err
}

func (rc RedisClient) connectInit() redis.Conn {
	conn := rc.pool.Get()
	return conn
}
func (rc RedisClient) Ping() bool {
	conn := rc.connectInit()
	defer conn.Close()
	pong, err := conn.Do("PING")
	if err != nil {
		// fmt.Println(err)
		return false
	}
	return pong == "PONG"
}
func (rc RedisClient) DisConnect() {
	rc.pool.Close()
}

func CreateRedisPool(addr string, maxActive, maxIdle int, wait bool) (pool *redis.Pool) {
	msgRedisConfig := addr
	pool = &redis.Pool{
		MaxActive:   maxActive,
		MaxIdle:     maxIdle,
		Wait:        wait,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", msgRedisConfig)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	return
}

func InitRedisClient(addr string, MaxActive, MaxIdle int, Wait bool) (Redis, error) {
	pool := CreateRedisPool(addr, MaxActive, MaxIdle, Wait)
	return RedisClient{pool, addr}, nil
}

func InitDefaultClient(addr string) (Redis, error) {
	return InitRedisClient(addr, 0, 3, true)
}
