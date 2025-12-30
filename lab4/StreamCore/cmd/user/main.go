package main

import (
	"StreamCore/internal/pkg/base"
	"StreamCore/internal/user"
	"StreamCore/kitex_gen/user/userservice"
	"log"
)

func main() {
	infra := base.GetInfraSet(
		base.WithDB(),
		base.WithCache())

	svr := userservice.NewServer(user.NewUserHandler(infra))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
