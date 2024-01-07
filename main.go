package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	discordInteractionAuth := newDiscordInteractionAuth()

	engine := gin.Default()

	engine.POST("/webhooks/discord/interactions", func(ctx *gin.Context) {
		if verified, err := discordInteractionAuth.verify(ctx); err != nil {
			ctx.String(http.StatusBadRequest, "failed to parse credentials: %v", err)

			return
		} else if !verified {
			ctx.String(http.StatusUnauthorized, "access denied")

			return
		}

		var request discordInteractionRequest
		if err := ctx.ShouldBindJSON(&request); err != nil {
			ctx.String(http.StatusBadRequest, "failed to parse request body: %v", err)

			return
		}

		switch request.Type {
		case discordInteractionRequestPing:
			ctx.JSON(http.StatusOK, &discordInteractionResponse{
				Type: discordInteractionResponsePong,
				Data: nil,
			})
		default:
			ctx.String(http.StatusNotImplemented, "interaction type %v not implemented", request.Type)
		}
	})

	if err := engine.Run(); err != nil {
		log.Fatalf("failed during server main loop: %v", err)
	}
}
