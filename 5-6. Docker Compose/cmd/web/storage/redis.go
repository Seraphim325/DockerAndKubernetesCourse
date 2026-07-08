package storage

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Addr        string
	Password    string
	Username    string
	DB          int
	MaxRetries  int
	DialTimeout time.Duration
	Timeout     time.Duration
}

func NewClient(c context.Context, conf Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         conf.Addr,
		Password:     conf.Password,
		Username:     conf.Username,
		DB:           conf.DB,
		MaxRetries:   conf.MaxRetries,
		DialTimeout:  conf.DialTimeout,
		WriteTimeout: conf.Timeout,
		ReadTimeout:  conf.Timeout,
	})

	if err := client.Ping(c).Err(); err != nil {
		log.Fatalf("Failed to establish connection with redis database: %s", err)
		return nil, err
	}

	return client, nil
}
