package main

import (
	"bytes"
	"context"
	"encoding/gob"
	"sync"
	"time"

	"github.com/go-redis/cache/v9"
	"github.com/redis/go-redis/v9"
)

var ErrCacheMiss = cache.ErrCacheMiss

type Cache interface {
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) bool
	Get(ctx context.Context, key string, value interface{}) error
	Set(ctx context.Context, key string, value interface{}, ttls ...time.Duration) error
}

type RedisCache struct {
	*cache.Cache
}

func NewRedisCache(ctx context.Context, addr string) (Cache, error) {
	client, err := newRedis(ctx, addr)
	if err != nil {
		return nil, err
	}

	return &RedisCache{
		Cache: cache.New(&cache.Options{
			Redis:      client,
			LocalCache: cache.NewTinyLFU(1000, time.Minute),
		}),
	}, nil
}

func newRedis(ctx context.Context, addr string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return client, nil
}

func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, ttls ...time.Duration) error {
	ttl := 1 * time.Hour
	if len(ttls) > 0 {
		ttl = ttls[0]
	}

	return c.Cache.Set(&cache.Item{
		Ctx:   ctx,
		Key:   key,
		Value: value,
		TTL:   ttl,
	})
}

type InMemoryCache struct {
	cache map[string][]byte // [key]value

	mu sync.RWMutex
}

func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{
		cache: make(map[string][]byte),
		mu:    sync.RWMutex{},
	}
}

func (c *InMemoryCache) Delete(ctx context.Context, key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.cache, key)

	return nil
}

func (c *InMemoryCache) Exists(ctx context.Context, key string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	_, ok := c.cache[key]

	return ok
}

func (c *InMemoryCache) Get(ctx context.Context, key string, value interface{}) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	v, ok := c.cache[key]
	if !ok {
		return ErrCacheMiss
	}

	return gob.NewDecoder(bytes.NewBuffer(v)).Decode(value)
}

func (c *InMemoryCache) Set(ctx context.Context, key string, value interface{}, ttls ...time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	buf := bytes.NewBuffer(nil)

	err := gob.NewEncoder(buf).Encode(value)
	if err != nil {
		return err
	}

	c.cache[key] = buf.Bytes()

	ttl := 1 * time.Hour
	if len(ttls) > 0 {
		ttl = ttls[0]
	}

	if ttl == 0 {
		ttl = 1 * time.Hour
	}

	time.AfterFunc(ttl, func() {
		c.Delete(ctx, key)
	})

	return nil
}
