package main

import (
	"fmt"
	"github.com/gauravbansal74/solution/config"
	"github.com/gauravbansal74/solution/route"
	"github.com/gauravbansal74/solution/server"
	"github.com/gauravbansal74/solution/shared/logger"
	"github.com/gauravbansal74/solution/shared/queue"
)

// Required Vars for Server
var (
	conf     config.Config
	redisapi *queue.RedisClient
)

// Main Function to start server
func main() {

	/// Load config from env, vars prefixed with APP_
	conf = config.LoadConfig()
	if conf.GoogleApiKey == "" {
		logger.Fatal("google-api-key", fmt.Errorf("Google API key is required"), "Google API key is required", nil)
	}

	// Load and Configure Logger for this Project
	logger.LoggerConfig(conf)

	// Init Redis Client
	redisapi = &queue.RedisClient{
		Host: conf.RedisHost,
		Port: conf.RedisPort,
		DB:   conf.RedisDB,
	}
	redisapi.Init()
	// redis test on server start
	redisServer := queue.ReadRedisClient()
	err := redisServer.Client.Set("onLoad", "Ok", 0).Err()
	if err != nil {
		// If we get any error in redis start throw log fatal and pass log
		logger.Fatal("redis-client", err, "Error while connecting to redis server", nil)
	} else {
		redisapi.Client.Del("onLoad")
		logger.Info("redis-client", "Redis-server connection tested successfully", nil)
	}
	// Init Redis MQ
	con := queue.Init(conf)
	// Go rout to read messages
	go con.GetPayload()

	h := route.LoadHandler()

	// start the Server
	server.StartServer(h, server.Server{
		ListenHost: conf.ListenHost,
		ListenPort: conf.ListenPort,
	})

}
