package main

import (
	"context"
	"log"

	"github.com/saurabh254/system-design-implementation/ratelimiter/internal/routes"
	"github.com/saurabh254/system-design-implementation/ratelimiter/internal/server"
	logger "github.com/saurabh254/system-design-implementation/ratelimiter/internal/utils"

	"github.com/saurabh254/system-design-implementation/ratelimiter/internal/config"

	_ "github.com/saurabh254/system-design-implementation/ratelimiter/docs"
)

// @title Rate Limiter API
// @version 1.0
// @description API for rate limiting service
// @host localhost:8080
// @BasePath /
func main() {
	logger.Init("INFO")
	handler := routes.NewRouter()
	config := config.Load()
	srv := server.New(config.GetServerAddress(), handler)

	if err := srv.Start(context.Background()); err != nil {
		log.Fatal(err)
	}
}
