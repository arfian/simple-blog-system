package main

import (
	"context"
	"time"

	"simple-blog-system/cmd/rest"
	"simple-blog-system/config"
	appSetup "simple-blog-system/internal/setup"
	"simple-blog-system/pkg/log"
)

func main() {
	// config init
	log.InitZeroLog()
	config.InitConfig()
	// conf := config.GetConfig()

	_, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	// app setup init
	setup := appSetup.Init()

	rest.StartServer(setup)
}
