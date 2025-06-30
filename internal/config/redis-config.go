package config

import (
	"context"
	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()

var Client = redis.NewClient(&redis.Options{
	Addr: "host.docker.internal.:6379",
})
