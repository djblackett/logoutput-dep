package main

import (
	"bufio"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	var str = uuid.NewString()

	r := gin.Default()
	url := os.Getenv("PING_PONG_SERVICE")
	message := os.Getenv("MESSAGE")

	if url == "" {
		log.Fatal("PING_PONG_SERVICE environment variable not set")
	}

	text, err := readConfig()
	fmt.Println(text, err)

	r.GET("/", func(c *gin.Context) {
		var timestamp = time.Now().Format(time.RFC3339)
		resp, err := http.Get(url + "/")

		if err != nil {
			log.Fatal(err)
		}

		defer resp.Body.Close()
		bodyBytes, err := io.ReadAll(resp.Body)
		c.String(http.StatusOK, "file content: "+text+"\nenv variable: MESSAGE="+message+"\n"+timestamp+" "+str+"\n"+"Ping / Pongs: "+string(bodyBytes))
	})

	r.GET("healthz", func(c *gin.Context) {
		_, err := http.Get(url + "/")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "cannot connect to Pingpong server"})
		} else {
			c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "connected to Pingpong server"})
		}
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r.Run()
}

func readConfig() (string, error) {
	var filename = "etc/config/informational.txt"

	file, err := os.Open(filename)

	if err != nil {
		fmt.Println("Error opening file:", err)
	}
	defer file.Close()

	var lines []string

	// Create a new scanner to read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	var MESSAGE = lines[0]

	// Check for errors during scanning
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}

	return MESSAGE, nil
}
