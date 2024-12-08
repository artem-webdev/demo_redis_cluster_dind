package redis_test_containers

import (
	"context"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type RepoRedisDemo struct {
	BaseSuite
}

func (st *RepoRedisDemo) Test_Universal() {
	type args struct {
		ctx        context.Context
		key        string
		value      interface{}
		expiration time.Duration
	}
	testsCases := []struct {
		name    string
		args    args
		wontErr bool
	}{
		{
			name: "set: проверка записи с универсальным клиентом",
			args: args{
				ctx:        context.Background(),
				key:        "key1",
				value:      "value1",
				expiration: 10 * time.Second,
			},
			wontErr: false,
		},
	}
	for _, tt := range testsCases {
		st.Suite.T().Run(tt.name, func(t *testing.T) {
			client := st.repoRedis.UniversalClient()
			err := client.Set(tt.args.ctx, tt.args.key, tt.args.value, tt.args.expiration).Err()
			if tt.wontErr {
				st.Suite.Require().Error(err)
			} else {
				st.Suite.Require().NoError(err)
			}
		})
	}
}

func (st *RepoRedisDemo) Test_Cluster() {
	type args struct {
		ctx        context.Context
		key        string
		value      interface{}
		expiration time.Duration
	}
	testsCases := []struct {
		name    string
		args    args
		wontErr bool
	}{
		{
			name: "set: проверка записи с кластерным клиентом клиентом",
			args: args{
				ctx:        context.Background(),
				key:        "key1",
				value:      "value1",
				expiration: 10 * time.Second,
			},
			wontErr: false,
		},
	}
	for _, tt := range testsCases {
		st.Suite.T().Run(tt.name, func(t *testing.T) {
			client := st.repoRedis.ClusterClient()
			err := client.Set(tt.args.ctx, tt.args.key, tt.args.value, tt.args.expiration).Err()
			if tt.wontErr {
				st.Suite.Require().Error(err)
			} else {
				st.Suite.Require().NoError(err)
			}
		})
	}
}

func TestRepoRedisDemoSuite(t *testing.T) {
	suite.Run(t, new(RepoRedisDemo))
}
