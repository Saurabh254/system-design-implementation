package main

import (
	"context"
	"log"

	"github.com/saurabh254/system-design-implementation/ratelimiter/internal/routes"
	"github.com/saurabh254/system-design-implementation/ratelimiter/internal/server"
	logger "github.com/saurabh254/system-design-implementation/ratelimiter/internal/utils"

	"github.com/saurabh254/system-design-implementation/ratelimiter/internal/config"
)

func main() {
	logger.Init("INFO")
	handler := routes.NewRouter()
	config := config.Load()
	srv := server.New(config.GetServerAddress(), handler)

	if err := srv.Start(context.Background()); err != nil {
		log.Fatal(err)
	}
}
