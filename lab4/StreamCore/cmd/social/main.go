package main

import (
	"StreamCore/internal/pkg/base"
	"StreamCore/internal/social"
	"StreamCore/kitex_gen/social/socialservice"
	"log"
)

func main() {
	infra := base.GetInfraSet(
		base.WithDB(),
		base.WithCache(),
		base.WithES(),
		base.WithMQ())

	svr := socialservice.NewServer(social.NewSocialHandler(infra))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
