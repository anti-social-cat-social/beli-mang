package main

import (
	"fmt"
	"halosuster/config"
	"halosuster/server"
	"log"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Load ENV from OS env or from .env variable
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	db := config.InitDb()

	fmt.Println(db.Ping())

	// Logger
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	r := gin.Default()

	// Initialize all routes
	server.NewRoute(r, db)

	// Start the server
	r.Run("0.0.0.0:8088")
}
