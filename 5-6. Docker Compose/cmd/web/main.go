package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"main/cmd/web/storage"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func main() {
	loadConfig()
	client := establishConnectionWithRedis()
	port := os.Getenv("SERVER_PORT")

	log.Printf("Server running on port %s\n", port)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		render(w, client)
	})

	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		log.Panic(err)
	}
}

func loadConfig() {
	if err := godotenv.Load(".env"); err != nil {
		log.Panic(err)
	}
}

func establishConnectionWithRedis() *redis.Client {
	conf := storage.Config{
		Addr:        fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Username:    os.Getenv("REDIS_USERNAME"),
		Password:    os.Getenv("REDIS_PASSWORD"),
		DB:          0,
		MaxRetries:  5,
		DialTimeout: 10 * time.Second,
		Timeout:     5 * time.Second,
	}

	client, err := storage.NewClient(context.Background(), conf)

	if err != nil {
		panic(err)
	}

	return client
}

func render(w http.ResponseWriter, client *redis.Client) {
	visits, err := client.Get(context.Background(), "visits").Result()

	if err == redis.Nil {
		visits = "0"
	} else if err != nil {
		log.Fatalf("Failed to get data: %s\n", err)
	}

	t, err := template.New("").Parse(fmt.Sprintf("Number of visits: %s", visits))

	if err = t.Execute(w, nil); err != nil {
		log.Print(err)
	}

	visitsNumber, err := strconv.Atoi(visits)

	if err != nil {
		log.Fatalf("Failed to convert visits: %s\n", err)
		client.Set(context.Background(), "visits", visits, 24*time.Hour)
	} else {
		client.Set(context.Background(), "visits", visitsNumber+1, 24*time.Hour)
	}
}
