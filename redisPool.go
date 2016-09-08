package redisPool

import (
	"time"

	"github.com/garyburd/redigo/redis"
)

type RedisPool struct {
	pool *redis.Pool
}

func NewRedisPool(server, password string) (*RedisPool, error) {

	if server == "" {
		server = ":6379"
	}

	pool := &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}

			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	return &RedisPool{pool}, nil
}

/*
删除数据
*/
func (r *RedisPool) DEL(key string) (int64, error) {
	c := r.pool.Get()
	defer c.Close()
	re, err := c.Do("DEL", key)
	if err != nil {
		return 0, err
	}
	return re.(int64), nil
}

/*
获得数据
*/
func (r *RedisPool) GET(key string) (interface{}, error) {
	c := r.pool.Get()
	defer c.Close()
	data, err := c.Do("GET", key)
	if err == nil {
		return data, nil
	}
	return nil, err
}

func (r *RedisPool) SETT(key string, value interface{}, express int64) bool {
	c := r.pool.Get()
	defer c.Close()
	_, err := c.Do("SET", key, value)
	if err == nil {
		_, err := c.Do("EXPIRE", key, express)
		if err == nil {
			return true
		}
	}
	return false
}
func (r *RedisPool) SET(key string, value interface{}) bool {
	c := r.pool.Get()
	defer c.Close()
	_, err := c.Do("SET", key, value)
	if err == nil {
		return true
	}
	return false
}
func (r *RedisPool) TTL(key string) int64 {
	c := r.pool.Get()
	defer c.Close()
	d, err := c.Do("TTL", key)
	if err == nil {
		return d.(int64)
	}
	return 0
}
