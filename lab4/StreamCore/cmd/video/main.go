package main

import (
	"StreamCore/internal/pkg/base"
	"StreamCore/internal/video"
	"StreamCore/kitex_gen/video/videoservice"
	"log"
)

func main() {
	infra := base.GetInfraSet(
		base.WithDB(),
		base.WithCache(),
		base.WithES(),
		base.WithMQ())

	svr := videoservice.NewServer(video.NewVideoHandler(infra))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
