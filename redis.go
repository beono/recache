package recache

import (
	"time"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

const tagPrefix = "tag"

// ErrKeyNotFound says that `key doesn't exist`
var ErrKeyNotFound = errors.New("key doesn't exist")

// RedisCache implements cache for redis
type RedisCache struct {
	redis *redis.Client
}

// NewRedisCache returns new RedisCache
func NewRedisCache(redis *redis.Client) RedisCache {
	return RedisCache{
		redis: redis,
	}
}

// Set saves data into redis
// Use ttl for `SETEX`-like behavior.
// Zero expiration means the key has no expiration time.
func (m *RedisCache) Set(id string, data interface{}, ttl uint, tags ...string) error {
	err := m.redis.Set(id, data, time.Duration(ttl)*time.Second).Err()
	if err != nil {
		return errors.Wrap(err, "can't set value")
	}
	if len(tags) == 0 {
		return nil
	}
	return errors.Wrap(m.setTags(id, tags...), "can't set tags")
}

// Get returns data from redis by key
func (m *RedisCache) Get(id string) ([]byte, error) {
	result, err := m.redis.Get(id).Bytes()
	if err == redis.Nil {
		return nil, ErrKeyNotFound
	}
	return result, errors.Wrap(err, "can't get value")
}

// ClearByTag deletes all keys that have the specified tag
func (m *RedisCache) ClearByTag(tag string) error {
	tag = m.prefixTag(tag)
	var cursor uint64
	for {
		keys, cursor, err := m.redis.SScan(tag, cursor, "", 0).Result()
		if err != nil {
			return errors.Wrap(err, "can't sscan")
		}
		pipe := m.redis.Pipeline()
		pipe.Del(tag)
		pipe.Del(keys...)
		if _, err = pipe.Exec(); err != nil {
			return errors.Wrap(err, "can't exec")
		}
		if cursor == 0 {
			return nil
		}
	}
}

func (m *RedisCache) setTags(id string, tags ...string) error {
	for _, tag := range tags {
		tag = m.prefixTag(tag)
		if err := m.redis.SAdd(tag, id).Err(); err != nil {
			return errors.Wrap(err, "can't add tag")
		}
	}
	return nil
}

func (m *RedisCache) prefixTag(tag string) string {
	return tagPrefix + ":" + tag
}
