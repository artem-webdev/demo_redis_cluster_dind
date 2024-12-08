package containers

import (
	"context"
	"errors"
	"fmt"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/redis"
	"io"
	"log"
	"strings"
	"time"
)

const (
	// RedisClusterUser юзер для теста default с остальными могут быть проблемы
	redisClusterUser = "default"
	// RedisClusterMinNode почему 3 , 3 это минимальное количество нод в кластере
	// для теста используется 3 мастера без реплик
	redisClusterMinNode = 3
	redisDefaultPort    = 6379
)

var (
	ErrExecComandInContainer = errors.New("error executing command in container")
)

type RedisClusterConfig struct {
	Addrs          []string
	Username       string
	Password       string
	PoolSize       int
	MaxActiveConns int
	PoolTimeout    time.Duration
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
}

// waitForClusterReady хелпер для ожидания готовности кластера
func waitForClusterReady(ctx context.Context, container *redis.RedisContainer, password string, timeout time.Duration) error {
	start := time.Now()
	for {
		clusterInfoCmd := []string{"redis-cli", "-a", password, "CLUSTER", "INFO"}
		errCode, stdout, err := container.Exec(ctx, clusterInfoCmd)
		if err != nil || errCode != 0 {
			return fmt.Errorf("failed to execute CLUSTER INFO: %w", err)
		}
		output, err := io.ReadAll(stdout)
		if err != nil {
			return fmt.Errorf("failed to read stdout: %w", err)
		}
		if strings.Contains(string(output), "cluster_state:ok") {
			return nil
		}
		if time.Since(start) > timeout {
			return fmt.Errorf("timeout waiting for cluster to be ready")
		}
		time.Sleep(1 * time.Second)
	}
}

func withRedisClusterCmdArgs(args ...string) testcontainers.CustomizeRequestOption {
	return func(req *testcontainers.GenericContainerRequest) error {
		if len(req.Cmd) == 0 {
			req.Cmd = []string{"redis-server"}
		}
		req.Cmd = append(req.Cmd, args...)
		return nil
	}
}

func SetupRedisCluster(ctx context.Context, password, image string) ([]*redis.RedisContainer, *RedisClusterConfig, error) {
	containers := make([]*redis.RedisContainer, redisClusterMinNode)
	externalAddrs := make([]string, redisClusterMinNode)
	internalAddrs := make([]string, redisClusterMinNode)
	for i := range redisClusterMinNode {
		container, err := redis.Run(ctx,
			image,
			withRedisClusterCmdArgs(
				"--cluster-enabled", "yes",
				"--cluster-config-file", "nodes.conf",
				"--cluster-node-timeout", "5000",
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
		ip, err := container.ContainerIP(ctx)
		if err != nil {
			return nil, nil, err
		}
		containers[i] = container
		externalAddrs[i] = strings.TrimPrefix(addr, "redis://")
		internalAddrs[i] = fmt.Sprintf("%s:%d", strings.TrimSpace(ip), redisDefaultPort)
	}
	// собираем кластер через redis-cli команду в первом контейнере
	clusterInitArgs := []string{"sh", "-c", fmt.Sprintf(
		"redis-cli -h %s -a %s --cluster create %s --cluster-replicas 0 --cluster-yes",
		internalAddrs[0], password, strings.Join(internalAddrs, " "),
	)}
	errCode, output, err := containers[0].Exec(ctx, clusterInitArgs)
	if err != nil || errCode != 0 {
		msg, errReader := io.ReadAll(output)
		if err != nil {
			return nil, nil, errReader
		}
		log.Println(string(msg))
		return nil, nil, ErrExecComandInContainer
	}
	// ждем пока кластер соберется
	if err = waitForClusterReady(ctx, containers[0], password, 30*time.Second); err != nil {
		log.Println(err.Error())
		return nil, nil, err
	}
	return containers, &RedisClusterConfig{
		Addrs:          externalAddrs,
		Username:       redisClusterUser,
		Password:       password,
		PoolSize:       10,
		MaxActiveConns: 10,
		PoolTimeout:    10 * time.Second,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}, nil
}
