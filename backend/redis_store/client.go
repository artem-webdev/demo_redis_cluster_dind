package redis_store

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

type Config struct {
	Addrs        []string
	Username     string
	Password     string
	PoolSize     int
	PoolTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type Repo struct {
	clusterClient   *redis.ClusterClient
	universalClient redis.UniversalClient
}

func New() *Repo {
	return &Repo{}
}

func (r *Repo) InitClusterClient(ctx context.Context, cnf *Config) error {
	clusterClient := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:        cnf.Addrs,
		Username:     cnf.Username,
		Password:     cnf.Password,
		PoolSize:     cnf.PoolSize,
		PoolTimeout:  cnf.PoolTimeout,
		ReadTimeout:  cnf.ReadTimeout,
		WriteTimeout: cnf.WriteTimeout,
	})
	_, err := clusterClient.Ping(ctx).Result()
	if err != nil {
		return err
	}
	log.Println("connected to Redis on ", "addr", cnf.Addrs)
	r.clusterClient = clusterClient
	return nil
}

func (r *Repo) ClusterClient() *redis.ClusterClient {
	return r.clusterClient
}

func (r *Repo) InitUniversalClient(ctx context.Context, cnf *Config) error {
	universalClient := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:        cnf.Addrs,
		Username:     cnf.Username,
		Password:     cnf.Password,
		PoolSize:     cnf.PoolSize,
		PoolTimeout:  cnf.PoolTimeout,
		ReadTimeout:  cnf.ReadTimeout,
		WriteTimeout: cnf.WriteTimeout,
	})
	_, err := universalClient.Ping(ctx).Result()
	if err != nil {
		return err
	}
	log.Println("connected to Redis on ", "addr", cnf.Addrs)
	r.universalClient = universalClient
	return nil
}

func (r *Repo) UniversalClient() redis.UniversalClient {
	return r.universalClient
}
