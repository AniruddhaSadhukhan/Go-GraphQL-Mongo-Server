package routes

import (
	"go-graphql-mongo-server/config"
	"go-graphql-mongo-server/logger"
	"strconv"
	"time"

	"github.com/sethvargo/go-limiter/httplimit"
	"github.com/sethvargo/go-limiter/memorystore"
)

var limiterMiddleware *httplimit.Middleware

func createLimiterMiddleware() {
	apiLimitPerSecond, err := strconv.ParseUint(config.ConfigManager.ApiLimitPerSecond, 10, 64)
	if err != nil || apiLimitPerSecond == 0 {
		logger.Log.Error("Error parsing apiLimitPerSecond: " + err.Error())
		return
	}

	limiterStore, err := memorystore.New(&memorystore.Config{
		// Number of API calls allowed per interval
		Tokens: apiLimitPerSecond,

		// Interval for which the limit is applied
		Interval: time.Second,
	})

	if err != nil {
		logger.Log.Error("Error creating limiter store: " + err.Error())
		return
	}

	limiterMiddleware, err = httplimit.NewMiddleware(limiterStore, httplimit.IPKeyFunc("X-Forwarded-For"))
	if err != nil {
		logger.Log.Error("Error creating limiter middleware: " + err.Error())
		return
	}
}
