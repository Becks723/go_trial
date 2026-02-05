package main

import (
	"StreamCore/config"
	"StreamCore/internal/interaction"
	"StreamCore/internal/pkg/base"
	"StreamCore/internal/pkg/constants"
	"StreamCore/kitex_gen/interaction/interactionservice"
	"StreamCore/pkg/util"
	"log"
	"net"

	"github.com/cloudwego/kitex/pkg/limit"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
)

var (
	infra       *base.InfraSet
	serviceName = constants.InteractionServiceName
	logPrefix   = "[interaction]"
)

func init() {
	config.Init(serviceName)
	infra = base.GetInfraSet(
		base.WithDB(),
		base.WithCache(),
		base.WithES(),
		base.WithMQ())
}

func main() {
	config := config.Instance()
	r, err := etcd.NewEtcdRegistry([]string{config.Etcd.Addr})
	if err != nil {
		log.Fatalf("%s NewEtcdRegistry error: %v", logPrefix, err)
	}
	listenAddr, ok := util.GetAvailablePort(config.Service.AddrList)
	if !ok {
		log.Fatalf("%s no port available", logPrefix)
	}
	addr, err := net.ResolveTCPAddr("tcp", listenAddr)
	if err != nil {
		log.Fatalf("%s ResolveTCPAddr error: %v", logPrefix, err)
	}

	svr := interactionservice.NewServer(
		interaction.NewInteractionHandler(infra),
		// 指定服务信息
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
			ServiceName: serviceName,
		}),
		// 开启多用复用
		server.WithMuxTransport(),
		// 指定RPC服务监听地址
		server.WithServiceAddr(addr),
		// 服务的注册与发现
		server.WithRegistry(r),
		// 设置限流
		server.WithLimit(&limit.Option{
			MaxConnections: constants.MaxConnections,
			MaxQPS:         constants.MaxQPS,
		}),
	)
	if err = svr.Run(); err != nil {
		log.Fatalf("%s server.Run error: %v", logPrefix, err)
	}
}
