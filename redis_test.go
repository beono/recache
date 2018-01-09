package recache

import (
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis"
)

var (
	cl  *redis.Client
	cache RedisCache
)

func TestMain(m *testing.M) {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	cl = redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	cache = NewRedisCache(cl)
	m.Run()
}

func TestRedisCache_Set(t *testing.T) {
	err := cache.Set("foo", "bar", 10)
	if err != nil {
		t.Errorf("unexpected error: %q", err)
	}

	res := cl.Get("foo").Val()
	if res != "bar" {
		t.Errorf("unexpected result: %q", res)
	}
}

func TestRedisCache_Get(t *testing.T) {
	cl.Set("foo", "bar", 0)

	res, err := cache.Get("foo")
	if err != nil {
		t.Errorf("unexpected error: %q", err)
	}

	if string(res) != "bar" {
		t.Errorf("unexpected result: %q", res)
	}
}

// TODO fix this test. Implement logic for that
//func TestRedisCache_Set_RemoveOldTags(t *testing.T) {
//	// we mark this cache entry with tag "tag1"
//	err := cache.Set("foo", "bar", 10, "tag1")
//	if err != nil {
//		t.Errorf("unexpected error: %q", err)
//	}
//
//	// then, we remark this cache entry with tag "tag2"
//	err = cache.Set("foo", "bar", 10, "tag2")
//	if err != nil {
//		t.Errorf("unexpected error: %q", err)
//	}
//
//	// clear by "tag1"
//	// it must not clear "foo" entry since "tag1" is no longer related to that entry.
//	err = cache.ClearByTag("tag1")
//	if err != nil {
//		t.Errorf("unexpected error: %q", err)
//	}
//
//	if res, _ := cache.Get("foo"); string(res) != "bar" {
//		t.Errorf("unexpected result: %q", string(res))
//	}
//}

func TestRedisCache_ClearByTag_ClearOne(t *testing.T) {

	// set up keys
	if err := cache.Set("foo", "bar", 10, "tag1", "tag2"); err != nil {
		t.Errorf("unexpected error: %q", err)
	}
	if res, _ := cache.Get("foo"); string(res) != "bar" {
		t.Errorf("unexpected result: %q", string(res))
	}

	if err := cache.Set("foo2", "bar", 10, "tag2", "tag3"); err != nil {
		t.Errorf("unexpected error: %q", err)
	}
	if res, _ := cache.Get("foo2"); string(res) != "bar" {
		t.Errorf("unexpected result: %q", string(res))
	}

	// clear by tag
	err := cache.ClearByTag("tag1")
	if err != nil {
		t.Errorf("unexpected error: %q", err)
	}

	// check that only ONE key has been flushed
	res, err := cache.Get("foo")
	if err != ErrKeyNotFound {
		t.Errorf("unexpected error: %q", err)
	}
	if string(res) != "" {
		t.Errorf("unexpected result: %q", string(res))
	}

	res, err = cache.Get("foo2")
	if err != nil {
		t.Errorf("unexpected error: %q", err)
	}
	if string(res) != "bar" {
		t.Errorf("unexpected result: %q", string(res))
	}
}

func TestRedisCache_ClearByTag_CleanAll(t *testing.T) {

	// set up keys
	if err := cache.Set("foo", "bar", 10, "tag1", "tag2"); err != nil {
		t.Errorf("unexpected error: %q", err)
	}
	if res, _ := cache.Get("foo"); string(res) != "bar" {
		t.Errorf("unexpected result: %q", string(res))
	}

	if err := cache.Set("foo2", "bar", 10, "tag2", "tag3"); err != nil {
		t.Errorf("unexpected error: %q", err)
	}
	if res, _ := cache.Get("foo2"); string(res) != "bar" {
		t.Errorf("unexpected result: %q", string(res))
	}

	// clear by tag
	err := cache.ClearByTag("tag2")
	if err != nil {
		t.Errorf("unexpected error: %q", err)
	}

	// check that ALL keys have been flushed
	_, err = cache.Get("foo")
	if err != ErrKeyNotFound {
		t.Errorf("unexpected error: %q", err)
	}

	_, err = cache.Get("foo2")
	if err != ErrKeyNotFound {
		t.Errorf("unexpected error: %q", err)
	}
}
