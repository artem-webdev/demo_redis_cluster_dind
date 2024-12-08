package containers

import (
	"context"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/redis"
	"strings"
	"time"
)

const (
	// RedisUser юзер для теста default с остальными могут быть проблемы
	redisUser = "default"
)

type RedisUniversalConfig struct {
	Addrs          []string
	Username       string
	Password       string
	PoolSize       int
	MaxActiveConns int
	PoolTimeout    time.Duration
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
}

func withRedisCmdArgs(args ...string) testcontainers.CustomizeRequestOption {
	return func(req *testcontainers.GenericContainerRequest) error {
		if len(req.Cmd) == 0 {
			req.Cmd = []string{"redis-server"}
		}
		req.Cmd = append(req.Cmd, args...)
		return nil
	}
}

func SetupRedis(ctx context.Context, password, image string) (*redis.RedisContainer, *RedisUniversalConfig, error) {
	container, err := redis.Run(ctx,
		image,
		withRedisCmdArgs(
			"--requirepass", password,
		),
		redis.WithLogLevel(redis.LogLevelVerbose),
	)
	if err != nil {
		return nil, nil, err
	}
	addr, err := container.ConnectionString(ctx)
	if err != nil {
		return nil, nil, err
	}
	return container, &RedisUniversalConfig{
		Addrs:          []string{strings.TrimPrefix(addr, "redis://")},
		Username:       redisUser,
		Password:       password,
		PoolSize:       10,
		MaxActiveConns: 10,
		PoolTimeout:    10 * time.Second,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}, nil
}
