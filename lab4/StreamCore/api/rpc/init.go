package rpc

import (
	"StreamCore/config"
	"StreamCore/internal/pkg/constants"
	"StreamCore/kitex_gen/interaction/interactionservice"
	"StreamCore/kitex_gen/social/socialservice"
	"StreamCore/kitex_gen/user/userservice"
	"StreamCore/kitex_gen/video/videoservice"
	"fmt"

	"github.com/cloudwego/kitex/client"
	etcd "github.com/kitex-contrib/registry-etcd"
)

var (
	userClient   userservice.Client
	videoClient  videoservice.Client
	iaClient     interactionservice.Client
	socialClient socialservice.Client
)

func Init() {
	initUserRPC()
	initVideoRPC()
	initSocialRPC()
	initInteractionRPC()
}

func initRPCClient[T any](serviceName string, newClientFunc func(string, ...client.Option) (T, error)) (*T, error) {
	config := config.Instance()
	r, err := etcd.NewEtcdResolver([]string{config.Etcd.Addr})
	if err != nil {
		return nil, fmt.Errorf("initRPCClient: error etcd.NewEtcdResolver: %w", err)
	}
	c, err := newClientFunc(serviceName,
		client.WithResolver(r),
		client.WithMuxConnection(constants.MuxConnection))
	if err != nil {
		return nil, fmt.Errorf("initRPCClient: error newClientFunc: %w", err)
	}
	return &c, nil
}
