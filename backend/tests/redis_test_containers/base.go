package redis_test_containers

import (
	"context"
	"github.com/artem-webdev/demo_redis_cluster_dind/redis_store"
	"github.com/artem-webdev/demo_redis_cluster_dind/tests/containers"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go/modules/redis"
)

const (
	imageRedis = "redis:7.2.4-alpine"
)

type BaseSuite struct {
	suite.Suite
	containerOneRedis      *redis.RedisContainer
	containersClusterRedis []*redis.RedisContainer
	repoRedis              *redis_store.Repo
}

func (st *BaseSuite) SetupSuite() {
	ctx := context.Background()

	containerRedis, cnfRd, err := containers.SetupRedis(ctx, "dev", imageRedis)
	st.Suite.Require().NoError(err)
	st.containerOneRedis = containerRedis
	redisCnf := &redis_store.Config{
		Addrs:        cnfRd.Addrs,
		Username:     cnfRd.Username,
		Password:     cnfRd.Password,
		PoolSize:     cnfRd.PoolSize,
		PoolTimeout:  cnfRd.PoolTimeout,
		ReadTimeout:  cnfRd.ReadTimeout,
		WriteTimeout: cnfRd.WriteTimeout,
	}
	st.repoRedis = redis_store.New()
	err = st.repoRedis.InitUniversalClient(ctx, redisCnf)
	st.Suite.Require().NoError(err)

	containersRedis, cnfCluster, err := containers.SetupRedisCluster(ctx, "dev", imageRedis)
	st.Suite.Require().NoError(err)
	st.containersClusterRedis = containersRedis
	cnfClusterRd := &redis_store.Config{
		Addrs:        cnfCluster.Addrs,
		Username:     cnfCluster.Username,
		Password:     cnfCluster.Password,
		PoolSize:     cnfCluster.PoolSize,
		PoolTimeout:  cnfCluster.PoolTimeout,
		ReadTimeout:  cnfCluster.ReadTimeout,
		WriteTimeout: cnfCluster.WriteTimeout,
	}
	err = st.repoRedis.InitClusterClient(ctx, cnfClusterRd)
	st.Suite.Require().NoError(err)
}

func (st *BaseSuite) TearDownSuite() {
	ctx := context.Background()
	err := st.containerOneRedis.Terminate(ctx)
	st.Suite.Require().NoError(err)
	for _, container := range st.containersClusterRedis {
		err = container.Terminate(ctx)
		st.Suite.Require().NoError(err)
	}
}
