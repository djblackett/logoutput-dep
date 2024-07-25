package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {

	redis_address := os.Getenv("REDIS_ADDRESS")

	client := redis.NewClient(&redis.Options{
		Addr:     redis_address + ":6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	ctx := context.Background()

	pings, err := client.Get(ctx, "pings").Result()
	var requestCount int

	if errors.Is(err, redis.Nil) {
		requestCount = 0
	} else if err != nil {
		panic(err)
	} else {
		requestCount, err = strconv.Atoi(pings)
		if err != nil {
			log.Fatalf("Error converting pings value to integer: %v", err)
		}
	}

	fmt.Println("pings:", requestCount)

	r := gin.Default()

	r.GET("/pingpong", func(c *gin.Context) {
		requestCount += 1
		c.String(http.StatusOK, strconv.Itoa(requestCount))

		err = client.Set(ctx, "pings", requestCount, 0).Err()
		if err != nil {
			panic(err)
		}

	})

	r.GET("/healthz", func(c *gin.Context) {
		err := client.Ping(ctx).Err()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "cannot connect to Redis"})
		} else {
			c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "connected to Redis"})
		}
	})

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, strconv.Itoa(requestCount))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = ":8000"
	}

	r.Run(port)
}
