package main

import (
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	engine := gin.Default()
	discordMount(engine)

	if err := engine.Run(); err != nil {
		log.Fatalf("failed during server main loop: %v", err)
	}
}
