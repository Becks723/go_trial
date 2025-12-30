package main

import (
	"StreamCore/internal/interaction"
	"StreamCore/internal/pkg/base"
	"StreamCore/kitex_gen/interaction/interactionservice"
	"log"
)

func main() {
	infra := base.GetInfraSet(
		base.WithDB(),
		base.WithCache())

	svr := interactionservice.NewServer(interaction.NewInteractionHandler(infra))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
