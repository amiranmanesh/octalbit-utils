package redis

import (
	"context"
	"errors"
	"time"

	"github.com/gomodule/redigo/redis"
)

// Redis is the connection pool holder
// and helper for calling methods
type Redis struct {
	Pool *redis.Pool
}

// InitRedis inits a Redis with connection pool
// use this Redis object to call its methods
func InitRedis(host, password string, maxIdle int, idleTimeout time.Duration) (*Redis, error) {

	p := &redis.Pool{
		MaxIdle:     maxIdle,
		IdleTimeout: idleTimeout,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", host,
				redis.DialPassword(password),
			)
		},
	}

	conn, err := p.Dial()
	if err != nil {
		return nil, err
	}
	if conn.Err() != nil {
		return nil, conn.Err()
	}
	defer func() { _ = conn.Close() }()

	return &Redis{Pool: p}, nil
}

// Set is for setting different kinds of values.
// If context expired before the redis job,
// then func returns immediately with context error
func (p *Redis) Set(ctx context.Context, key string, value interface{}) error {

	c, ctxErr := p.Pool.GetContext(ctx)
	if ctxErr != nil {
		return ctxErr
	}
	defer func() { _ = c.Close() }()

	_, err := c.Do("SET", key, value)
	return err
}

// SetWithTTl is for setting value with TTL bind to it.
// If context expired before the redis job,
// then func returns immediately with context error
func (p *Redis) SetWithTTl(ctx context.Context, key string, value interface{}, ttl time.Duration) error {

	c, ctxErr := p.Pool.GetContext(ctx)
	if ctxErr != nil {
		return ctxErr
	}
	defer func() { _ = c.Close() }()

	_, err := c.Do("SETEX", key, uint(ttl.Seconds()), value)
	return err
}

// Incr is for increasing a key's value.
// If context expired before the redis job,
// then func returns immediately with context error
func (p *Redis) Incr(ctx context.Context, key string) (int, error) {

	c, ctxErr := p.Pool.GetContext(ctx)
	if ctxErr != nil {
		return 0, ctxErr
	}
	defer func() { _ = c.Close() }()

	return redis.Int(c.Do("INCR", key))
}

// IncrWithTTl increases a key's value and also updates its TTL.
// If context expired before the redis job,
// then func returns immediately with context error
func (p *Redis) IncrWithTTl(ctx context.Context, key string, ttl time.Duration) (int, error) {

	c, ctxErr := p.Pool.GetContext(ctx)
	if ctxErr != nil {
		return 0, ctxErr
	}
	defer func() {
		_, _ = c.Do("EXPIRE", key, uint(ttl.Seconds()))
		_ = c.Close()
	}()

	return redis.Int(c.Do("INCR", key))
}

// SetTTl is for setting expiration ttl to a key.
// If context expired before the redis job,
// then func returns immediately with context error
func (p *Redis) SetTTl(ctx context.Context, key string, ttl time.Duration) error {

	c, ctxErr := p.Pool.GetContext(ctx)
	if ctxErr != nil {
		return ctxErr
	}
	defer func() { _ = c.Close() }()

	_, err := c.Do("EXPIRE", key, uint(ttl.Seconds()))
	return err
}

// GetBytes gets value as bytes.
// If context expired before the redis job,
// then func returns immediately with context error
func (p *Redis) GetBytes(ctx context.Context, key string) ([]byte, error) {

	c, ctxErr := p.Pool.GetContext(ctx)
	if ctxErr != nil {
		return nil, ctxErr
	}
	defer func() { _ = c.Close() }()

	return redis.Bytes(c.Do("GET", key))
}

// GetString gets value as string.
// If context expired before the redis job,
// then func returns immediately with context error
func (p *Redis) GetString(ctx context.Context, key string) (string, error) {

	c, ctxErr := p.Pool.GetContext(ctx)
	if ctxErr != nil {
		return "", ctxErr
	}
	defer func() { _ = c.Close() }()

	return redis.String(c.Do("GET", key))
}

// GetInt gets value as int.
// If context expired before the redis job,
// then func returns immediately with context error
func (p *Redis) GetInt(ctx context.Context, key string) (int, error) {

	c, ctxErr := p.Pool.GetContext(ctx)
	if ctxErr != nil {
		return 0, ctxErr
	}
	defer func() { _ = c.Close() }()

	return redis.Int(c.Do("GET", key))
}

// Delete is for deleting one or several keys
// deletes all the existence keys.
// If context expired before the redis job,
// then func returns immediately with context error
func (p *Redis) Delete(ctx context.Context, keys ...string) error {

	c, ctxErr := p.Pool.GetContext(ctx)
	if ctxErr != nil {
		return ctxErr
	}
	defer func() { _ = c.Close() }()
	isAllDeleted := true

	for _, key := range keys {
		resp, err := c.Do("DEL", key)
		if err != nil {
			return err
		}
		if resp == 0 {
			isAllDeleted = false
		}
	}

	if !isAllDeleted {
		return errors.New("some keys does not exist")
	}
	return nil
}

// DeleteByPattern is for deleting keys by pattern.
// If context expired before the redis job,
// then func returns immediately with context error
func (p *Redis) DeleteByPattern(ctx context.Context, pattern string) error {

	c, ctxErr := p.Pool.GetContext(ctx)
	if ctxErr != nil {
		return ctxErr
	}
	defer func() { _ = c.Close() }()

	b, err := c.Do("KEYS", pattern)
	if err != nil {
		return err
	}

	keys := b.([]interface{})
	for _, key := range keys {
		_, err = c.Do("DEL", string(key.([]byte)))
		if err != nil {
			return err
		}
	}

	return nil
}

// AddToSet is for adding a value to a set
// returns true if value was not in the set before.
// If context expired before the redis job,
// then func returns immediately with context error
func (p *Redis) AddToSet(ctx context.Context, setName string, value string) (int, error) {

	c, ctxErr := p.Pool.GetContext(ctx)
	if ctxErr != nil {
		return 0, ctxErr
	}
	defer func() { _ = c.Close() }()

	_, err := c.Do("SADD", setName, value)
	if err != nil {
		return 0, err
	}

	return redis.Int(c.Do("SCARD", setName))
}

// AddToSetWithTTl is for adding a value to a set with ttl
// returns true if value was not in the set before.
// If context expired before the redis job,
// then func returns immediately with context error
func (p *Redis) AddToSetWithTTl(ctx context.Context, setName string, value string, ttl time.Duration) (int, error) {

	c, ctxErr := p.Pool.GetContext(ctx)
	if ctxErr != nil {
		return 0, ctxErr
	}
	defer func() {
		_, _ = c.Do("EXPIRE", setName, uint(ttl.Seconds()))
		_ = c.Close()
	}()

	_, err := c.Do("SADD", setName, value)
	if err != nil {
		return 0, err
	}

	return redis.Int(c.Do("SCARD", setName))
}
