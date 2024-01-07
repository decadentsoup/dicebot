package main

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

type discordInteractionRequest struct {
	Type int `json:"type"`
	Data any `json:"data,omitempty"`
}

const (
	discordInteractionRequestPing                           = 1
	discordInteractionRequestApplicationCommand             = 2
	discordInteractionRequestMessageComponent               = 3
	discordInteractionRequestApplicationCommandAutocomplete = 4
	discordInteractionRequestModalSubmit                    = 5
)

type discordInteractionResponse struct {
	Type int `json:"type"`
	Data any `json:"data,omitempty"`
}

const (
	discordInteractionResponsePong                                 = 1
	discordInteractionResponseChannelMessageWithSource             = 4
	discordInteractionResponseDeferredChannelMessageWithSource     = 5
	discordInteractionResponseDeferredUpdateMessage                = 6
	discordInteractionResponseUpdateMessage                        = 7
	discordInteractionResponseApplicationCommandAutocompleteResult = 8
	discordInteractionResponseModal                                = 9
	discordInteractionResponsePremiumRequired                      = 10
)

type discordInteractionAuth struct {
	applicationPublicKey ed25519.PublicKey
}

func newDiscordInteractionAuth() discordInteractionAuth {
	applicationPublicKey := make([]byte, ed25519.PublicKeySize)
	if _, err := hex.Decode(applicationPublicKey, []byte(os.Getenv("DISCORD_APPLICATION_PUBLIC_KEY"))); err != nil {
		log.Fatalf("failed to parse environment variable DISCORD_APPLICATION_PUBLIC_KEY: %v", err)
	}

	return discordInteractionAuth{
		applicationPublicKey: applicationPublicKey,
	}
}

func (auth *discordInteractionAuth) verify(ctx *gin.Context) (bool, error) {
	signature := make([]byte, ed25519.SignatureSize)
	if _, err := hex.Decode(signature, []byte(ctx.GetHeader("X-Signature-Ed25519"))); err != nil {
		return false, fmt.Errorf("failed to parse header X-Signature-Ed25519: %w", err)
	}

	timestamp := ctx.GetHeader("X-Signature-Timestamp")

	body, err := peekBody(ctx)
	if err != nil {
		return false, err
	}

	return ed25519.Verify(auth.applicationPublicKey, []byte(timestamp+string(body)), signature), nil
}
