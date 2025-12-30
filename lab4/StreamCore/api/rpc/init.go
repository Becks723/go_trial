package rpc

import (
	"StreamCore/internal/pkg/constants"
	"StreamCore/kitex_gen/interaction/interactionservice"
	"StreamCore/kitex_gen/social/socialservice"
	"StreamCore/kitex_gen/user/userservice"
	"StreamCore/kitex_gen/video/videoservice"
	"StreamCore/pkg/env"
	"errors"
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

const (
	MuxConnection = 1

	UserServiceName        = "user"
	VideoServiceName       = "video"
	InteractionServiceName = "interaction"
	SocialServiceName      = "social"
)

func initRPCClient[T any](serviceName string, newClientFunc func(string, ...client.Option) (T, error)) (*T, error) {
	addr := env.Instance().Etcd_Addr
	if addr == "" {
		return nil, errors.New("env etcd addr null")
	}
	r, err := etcd.NewEtcdResolver([]string{addr})
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
